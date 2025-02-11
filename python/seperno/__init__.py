import sys
import ctypes
import os

# Determine the shared library extension based on the OS
system = sys.platform
if system == "linux":
    lib_ext = "so"
elif system == "darwin":
    lib_ext = "dylib"
elif system == "win32":
    lib_ext = "dll"
else:
    raise OSError(f"Unsupported platform: {system}")

_lib_name = f"seperno.{lib_ext}"
_lib_path = os.path.join(os.path.dirname(__file__), _lib_name)
seperno = ctypes.CDLL(_lib_path)