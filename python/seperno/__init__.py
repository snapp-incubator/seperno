import sys
import ctypes
import os

# Determine shared library extension based on OS
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

# Define argument types for NormalizeText
seperno.NormalizeText.argtypes = [
    ctypes.c_char_p,
    ctypes.c_bool, ctypes.c_bool, ctypes.c_bool,
    ctypes.c_bool, ctypes.c_bool, ctypes.c_bool, ctypes.c_bool
]
seperno.NormalizeText.restype = ctypes.c_char_p

def normalize_text(
    text,
    convert_half_space=False,
    combine_space=False,
    remove_outer_space=False,
    remove_url=False,
    normalize_punctuations=False,
    end_with_eol=False,
    int_to_word=False
):
    """Python wrapper for the Go-based text normalizer."""
    return seperno.NormalizeText(
        text.encode("utf-8"),
        ctypes.c_bool(convert_half_space),
        ctypes.c_bool(combine_space),
        ctypes.c_bool(remove_outer_space),
        ctypes.c_bool(remove_url),
        ctypes.c_bool(normalize_punctuations),
        ctypes.c_bool(end_with_eol),
        ctypes.c_bool(int_to_word),
    ).decode("utf-8")