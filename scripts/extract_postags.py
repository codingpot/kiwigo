#!/usr/bin/env python3
"""
Extract POS tags from Kiwi C++ source code using tree-sitter.

This script parses the Kiwi C++ sources to extract the POSTag enum and the
tag<->string conversion tables, then generates a Go file with the corresponding
POS type constants.

The mapping is derived rather than guessed: enum values are evaluated the way
the C++ compiler evaluates them (including aliases such as `pv = p` and
`pa = p + 1`), and the resulting numeric value is fed through a Python
transcription of `tagToString`. Tags that only `tagRToString` can produce
(the regular-conjugation `-R` variants) are parsed from that function.

Usage:
    python scripts/extract_postags.py [version]

Example:
    python scripts/extract_postags.py v0.23.2
"""

import re
import sys
import urllib.request
from pathlib import Path

import tree_sitter_cpp as tscpp
from tree_sitter import Language, Parser

# Initialize parser
CPP_LANGUAGE = Language(tscpp.language())

# GitHub raw content URL template
KIWI_TYPES_H_URL = "https://raw.githubusercontent.com/bab2min/Kiwi/{version}/include/kiwi/Types.h"
KIWI_UTILS_CPP_URL = "https://raw.githubusercontent.com/bab2min/Kiwi/{version}/src/Utils.cpp"
KIWI_STRUTILS_H_URL = "https://raw.githubusercontent.com/bab2min/Kiwi/{version}/src/StrUtils.h"

# Enumerators that are not POS tags themselves.
SKIPPED_ENUMERATORS = {
    "max",  # size marker
    "irregular",  # bit flag
    "unknown_feat_ha",  # internal use
}

# Bases whose `<base>i` spelling maps to a POS_<BASE>_I constant.
IRREGULAR_BASES = ("vv", "va", "vx", "xsa", "pv", "pa")

# Go constant names for tag strings that toPOSTag accepts but that are not
# spelled like an identifier.
PUNCTUATION_TAG_NAMES = {
    "^": "POS_CARET",
}

# Backwards-compatible aliases kept for downstream code. These cannot be derived
# from the C++ source; they exist because kiwigo used to spell them this way.
COMPAT_ALIASES = [
    ("POS_USER_0", "POS_USER0"),
    ("POS_USER_1", "POS_USER1"),
    ("POS_USER_2", "POS_USER2"),
    ("POS_USER_3", "POS_USER3"),
    ("POS_USER_4", "POS_USER4"),
]


def fetch_file(url: str) -> str:
    """Fetch file content from URL."""
    with urllib.request.urlopen(url) as response:
        return response.read().decode("utf-8")


def extract_enum_values(source: str, enum_name: str) -> list[dict]:
    """Extract enum values from C++ source code using tree-sitter."""
    parser = Parser(CPP_LANGUAGE)
    tree = parser.parse(bytes(source, "utf8"))

    results = []

    def traverse(node):
        if node.type == "enum_specifier":
            # Check if this is the enum we're looking for
            name_node = node.child_by_field_name("name")
            if name_node and name_node.text.decode() == enum_name:
                # Find the body
                body = node.child_by_field_name("body")
                if body:
                    for child in body.children:
                        if child.type == "enumerator":
                            name = child.child_by_field_name("name")
                            value = child.child_by_field_name("value")
                            if name:
                                entry = {
                                    "name": name.text.decode(),
                                    "value": value.text.decode() if value else None,
                                }
                                results.append(entry)
        for child in node.children:
            traverse(child)

    traverse(tree.root_node)
    return results


def resolve_enum_values(tags: list[dict]) -> dict[str, int]:
    """Evaluate C++ enumerator initializers to their numeric values.

    Follows C++ rules: an enumerator without an initializer is the previous
    value plus one, and an initializer may reference earlier enumerators.
    """
    resolved: dict[str, int] = {}
    previous = -1

    for tag in tags:
        expr = tag["value"]
        if expr is None:
            value = previous + 1
        else:
            # Initializers in POSTag only use earlier enumerators, integer
            # literals and the `+` / `|` operators.
            if not re.fullmatch(r"[\w\s+|()x0-9A-Fa-f]+", expr):
                raise ValueError(f"unsupported enumerator initializer: {expr!r}")
            value = eval(expr, {"__builtins__": {}}, dict(resolved))  # noqa: S307

        resolved[tag["name"]] = value
        previous = value

    return resolved


def extract_function_body(source: str, signature: str) -> str:
    """Return the brace-delimited body of the first function matching `signature`."""
    match = re.search(signature, source)
    if not match:
        return ""

    start = source.find("{", match.end())
    if start == -1:
        return ""

    depth = 0
    for i in range(start, len(source)):
        if source[i] == "{":
            depth += 1
        elif source[i] == "}":
            depth -= 1
            if depth == 0:
                return source[start : i + 1]
    return ""


def extract_tag_strings(body: str) -> list[str]:
    """Extract the `tags[]` string table from a tagToString-like function body."""
    pattern = r"static\s+const\s+char\s*\*\s*tags\s*\[\]\s*=\s*\{([^}]+)\}"
    match = re.search(pattern, body, re.DOTALL)

    if not match:
        return []

    return re.findall(r'"([^"]+)"', match.group(1))


def extract_case_returns(body: str) -> dict[str, str]:
    """Extract `case POSTag::x: return "Y";` pairs from a function body."""
    pattern = r'case\s+POSTag::(\w+)\s*:\s*return\s+"([^"]*)"'
    return dict(re.findall(pattern, body))


def extract_accepted_tags(body: str) -> list[str]:
    """Extract the tag strings that `toPOSTag` accepts, in source order.

    toPOSTag is the upstream counterpart of ParsePOSType, so it defines which
    strings are valid input. Anything it does not recognise falls through to
    `return POSTag::max`, which is upstream's way of rejecting the string.
    """
    pattern = r'tagStr\s*==\s*u"([^"]*)"\s*\)\s*return'
    seen: list[str] = []
    for text in re.findall(pattern, body):
        if text not in seen:
            seen.append(text)
    return seen


def make_tag_to_string(tag_strings: list[str], irregular_cases: dict[str, str],
                       values: dict[str, int], irregular_flag: int, max_value: int):
    """Build a Python transcription of the C++ `tagToString`.

    Returns None where upstream has no tag string for the value. tagToString
    guards its table lookup with `assert(t < POSTag::max)` and its irregular
    branch falls through to a "@" sentinel, so those values are not tags that
    upstream ever converts; `toPOSTag` rejects "@" as well.
    """
    # Map the numeric value of each irregular case label to its return string.
    irregular_by_value = {values[name]: text for name, text in irregular_cases.items()}

    def tag_to_string(value: int) -> str | None:
        if value & irregular_flag:
            return irregular_by_value.get(value & ~irregular_flag)
        if value >= max_value:
            return None
        return tag_strings[value]

    return tag_to_string


def go_constant_name(cpp_name: str) -> str:
    """Map a C++ enumerator name to its Go constant name."""
    if cpp_name == "unknown":
        return "POS_UNKNOWN"

    if cpp_name.endswith("i") and cpp_name[:-1] in IRREGULAR_BASES:
        return f"POS_{cpp_name[:-1].upper()}_I"

    return f"POS_{cpp_name.upper()}"


def build_tag_entries(tags: list[dict], values: dict[str, int], tag_to_string,
                      regular_cases: dict[str, str],
                      accepted: list[str]) -> tuple[list[tuple[str, str]], list[tuple[str, str]]]:
    """Build the (go_name, go_value) lists for the generated constants.

    Returns the tags Kiwi can emit and, separately, the input-only aliases that
    upstream's toPOSTag accepts but never produces.
    """
    emitted: list[tuple[str, str]] = []
    seen_names: set[str] = set()
    seen_values: set[str] = set()

    def add(go_name: str, go_value: str, target: list[tuple[str, str]]) -> None:
        if go_name in seen_names:
            return
        seen_names.add(go_name)
        seen_values.add(go_value)
        target.append((go_name, go_value))

    for tag in tags:
        cpp_name = tag["name"]
        if cpp_name in SKIPPED_ENUMERATORS:
            continue

        go_value = tag_to_string(values[cpp_name])
        if go_value is None:
            continue

        add(go_constant_name(cpp_name), go_value, emitted)

    # `-R` variants are only reachable through tagRToString, so they are not
    # derivable from the enum alone.
    for cpp_name, text in regular_cases.items():
        add(f"POS_{cpp_name.upper()}_R", text, emitted)

    # Strings toPOSTag accepts as input but tagToString never returns, such as
    # "V" and "A" for POSTag::p or "NF"/"NV"/"NA"/"UNK" for POSTag::unknown.
    aliases: list[tuple[str, str]] = []
    for text in accepted:
        if text in seen_values:
            continue
        go_name = PUNCTUATION_TAG_NAMES.get(text, f"POS_{text.upper()}")
        add(go_name, text, aliases)

    return emitted, aliases


def generate_go_file(emitted: list[tuple[str, str]], aliases: list[tuple[str, str]],
                     version: str) -> str:
    """Generate Go source file with POS type definitions."""
    entries = emitted + aliases
    lines = []
    lines.append("// Code generated by scripts/extract_postags.py; DO NOT EDIT.")
    lines.append(f"// Source: Kiwi {version}")
    lines.append("//")
    lines.append("// To regenerate, run: make sync-postypes")
    lines.append("")
    lines.append("package kiwi")
    lines.append("")
    lines.append('import "fmt"')
    lines.append("")
    lines.append("type POSType string")
    lines.append("")
    lines.append("const (")

    # gofmt aligns each run of declarations separated by a blank line on its
    # own, so each group is padded to its own longest name.
    def declare_group(group: list[tuple[str, str]]) -> list[str]:
        width = max(len(name) for name, _ in group)
        return [
            f'\t{go_name}{" " * (width - len(go_name) + 1)}POSType = "{go_value}"'
            for go_name, go_value in group
        ]

    lines.extend(declare_group(emitted))
    lines.append("")
    lines.append("\t// Accepted by toPOSTag but never produced by Kiwi itself.")
    lines.extend(declare_group(aliases))

    lines.append(")")
    lines.append("")

    # Backwards-compatible aliases.
    lines.append("// Deprecated: these aliases are kept for backwards compatibility.")
    lines.append("const (")
    alias_len = max(len(old) for old, _ in COMPAT_ALIASES)
    for old_name, new_name in COMPAT_ALIASES:
        padding = " " * (alias_len - len(old_name) + 1)
        lines.append(f"\t{old_name}{padding}= {new_name}")
    lines.append(")")
    lines.append("")

    # Generate isValid function. This mirrors which strings upstream's toPOSTag
    # accepts, plus the strings tagToString/tagRToString can return. A Go switch
    # requires distinct case values, and several constants are aliases sharing
    # one string, so deduplicate by value.
    lines.append("func (p POSType) isValid() bool {")
    lines.append("\tswitch p {")
    lines.append("\tcase")

    valid_tags = []
    seen_values: set[str] = set()
    for go_name, go_value in entries:
        if go_value in seen_values:
            continue
        seen_values.add(go_value)
        valid_tags.append(go_name)

    lines.append(",\n".join(f"\t\t{tag}" for tag in valid_tags) + ":")
    lines.append("\t\treturn true")
    lines.append("\tdefault:")
    lines.append("\t\treturn false")
    lines.append("\t}")
    lines.append("}")
    lines.append("")

    # Generate ParsePOSType function
    lines.append("// ParsePOSType return POS Tag for result")
    lines.append("func ParsePOSType(t string) (POSType, error) {")
    lines.append("\tpos := POSType(t)")
    lines.append("\tif !pos.isValid() {")
    lines.append('\t\treturn POS_UNKNOWN, fmt.Errorf("POS type parse err. input type: %s", t)')
    lines.append("\t}")
    lines.append("\treturn pos, nil")
    lines.append("}")

    return "\n".join(lines) + "\n"


def main():
    version = sys.argv[1] if len(sys.argv) > 1 else "v0.23.2"

    print(f"Extracting POS tags from Kiwi {version}...")

    # Fetch source files
    types_url = KIWI_TYPES_H_URL.format(version=version)
    utils_url = KIWI_UTILS_CPP_URL.format(version=version)

    print(f"Fetching {types_url}")
    types_source = fetch_file(types_url)

    print(f"Fetching {utils_url}")
    utils_source = fetch_file(utils_url)

    strutils_url = KIWI_STRUTILS_H_URL.format(version=version)
    print(f"Fetching {strutils_url}")
    strutils_source = fetch_file(strutils_url)

    # Extract enum values
    print("Parsing POSTag enum...")
    tags = extract_enum_values(types_source, "POSTag")
    print(f"Found {len(tags)} enum values")

    values = resolve_enum_values(tags)
    if "irregular" not in values:
        raise SystemExit("POSTag::irregular not found; cannot decode irregular tags")

    # Extract the conversion tables.
    print("Parsing tag strings...")
    tag_to_string_body = extract_function_body(utils_source, r"const\s+char\s*\*\s*tagToString\s*\(")
    tag_strings = extract_tag_strings(tag_to_string_body)
    if not tag_strings:
        raise SystemExit("tagToString tag table not found")
    print(f"Found {len(tag_strings)} tag strings")

    irregular_cases = extract_case_returns(tag_to_string_body)

    tag_r_to_string_body = extract_function_body(utils_source, r"const\s+char\s*\*\s*tagRToString\s*\(")
    regular_cases = extract_case_returns(tag_r_to_string_body)
    print(f"Found {len(regular_cases)} regular-conjugation tags")

    # toPOSTag defines which strings upstream accepts as input, which is the
    # role ParsePOSType plays on the Go side.
    to_pos_tag_body = extract_function_body(strutils_source, r"inline\s+POSTag\s+toPOSTag\s*\(\s*std::u16string_view")
    accepted = extract_accepted_tags(to_pos_tag_body)
    if not accepted:
        raise SystemExit("toPOSTag not found; cannot determine accepted tag strings")
    print(f"Found {len(accepted)} accepted tag strings")

    tag_to_string = make_tag_to_string(
        tag_strings, irregular_cases, values, values["irregular"], values["max"]
    )
    emitted, aliases = build_tag_entries(tags, values, tag_to_string, regular_cases, accepted)

    # Generate Go file
    go_content = generate_go_file(emitted, aliases, version)

    # Write output
    output_path = Path("postype_generated.go")
    output_path.write_text(go_content)
    print(f"Generated {output_path}")

    # Print summary
    print("\nExtracted tags:")
    for tag in tags:
        cpp_name = tag["name"]
        text = None if cpp_name in SKIPPED_ENUMERATORS else tag_to_string(values[cpp_name])
        if text is None:
            print(f"  {cpp_name:20s} -> (no tag string; skipped)")
            continue
        print(f'  {cpp_name:20s} -> {go_constant_name(cpp_name):20s} = "{text}"')
    for cpp_name, text in regular_cases.items():
        print(f'  {cpp_name + " (R)":20s} -> {"POS_" + cpp_name.upper() + "_R":20s} = "{text}"')
    for go_name, text in aliases:
        print(f'  {"(toPOSTag only)":20s} -> {go_name:20s} = "{text}"')


if __name__ == "__main__":
    main()
