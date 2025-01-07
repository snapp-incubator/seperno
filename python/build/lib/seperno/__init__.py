import ctypes
import os

# Load the shared library
_lib_path = os.path.join(os.path.dirname(__file__), "seperno.so")
seperno = ctypes.CDLL(_lib_path)

# Set argument types for the NormalizeText function
seperno.NormalizeText.argtypes = [
    ctypes.c_char_p, ctypes.c_bool, ctypes.c_bool, ctypes.c_bool,
    ctypes.c_bool, ctypes.c_bool, ctypes.c_bool
]
seperno.NormalizeText.restype = ctypes.c_char_p

def normalize_text(text, convert_half_space=False, combine_space=False, remove_outer_space=False, remove_url=False, normalize_punctuations=False, end_with_eol=False):
    """Python wrapper for the Go-based text normalizer."""
    return seperno.NormalizeText(
        text.encode('utf-8'),
        ctypes.c_bool(convert_half_space),
        ctypes.c_bool(combine_space),
        ctypes.c_bool(remove_outer_space),
        ctypes.c_bool(remove_url),
        ctypes.c_bool(normalize_punctuations),
        ctypes.c_bool(end_with_eol)
    ).decode('utf-8')

# Example usage
if __name__ == "__main__":
    result = normalize_text("Hello  World!", convert_half_space=True)
    print("Normalized:", result)