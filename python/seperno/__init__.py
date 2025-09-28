import os
import platform
import ctypes

# Determine the shared library path based on the platform
system = platform.system()
if system == "Linux":
    lib_name = "seperno.so"
    lib_path = os.path.join(os.path.dirname(__file__), "linux", lib_name)
elif system == "Darwin":  # macOS
    lib_name = "seperno.dylib"
    lib_path = os.path.join(os.path.dirname(__file__), "macos", lib_name)
else:
    raise RuntimeError(f"Unsupported platform: {system}")

# Load the shared library
if not os.path.exists(lib_path):
    raise RuntimeError(f"Shared library not found: {lib_path}")

seperno = ctypes.CDLL(lib_path)

# Define argument types for NormalizeText
seperno.NormalizeText.argtypes = [
    ctypes.c_char_p,
    ctypes.c_bool, ctypes.c_bool, ctypes.c_bool,
    ctypes.c_bool, ctypes.c_bool, ctypes.c_bool,
    ctypes.c_bool, ctypes.c_char_p
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
    int_to_word=False,
    number_language="en"
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
        number_language.encode("utf-8"),
    ).decode("utf-8")