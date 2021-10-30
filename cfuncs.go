package kiwi

/*
#cgo LDFLAGS: -l kiwi

int KiwiReaderBridge(int lineNumber, char *buffer, void *userData) {
	int KiwiReaderImpl(int, char*, void*);
  return KiwiReaderImpl(lineNumber, buffer, userData);
}
*/
import "C"
