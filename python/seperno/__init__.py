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


seperno.CosineSimilarity.argtypes = [
    ctypes.POINTER(ctypes.c_char_p),
    ctypes.POINTER(ctypes.POINTER(ctypes.c_char_p)),
    ctypes.c_int, ctypes.c_int, ctypes.c_int
]
seperno.CosineSimilarity.restype = ctypes.POINTER(ctypes.c_double)


def cosine_similarity(text, list_text):
    # Convert text to a C array of strings
    text_c = (ctypes.c_char_p * len(text))(*[t.encode('utf-8') for t in text])

    # Convert list_text to a C array of C arrays of strings
    list_text_c = (ctypes.POINTER(ctypes.c_char_p) * len(list_text))()
    for i, lst in enumerate(list_text):
        inner_array = (ctypes.c_char_p * len(lst))(*[s.encode('utf-8') for s in lst])
        list_text_c[i] = inner_array

    # Call the Go function
    similarities_ptr = seperno.CosineSimilarity(
        text_c, list_text_c, len(text), len(list_text), len(list_text[0])
    )

    # Convert the C result array back to a Python list
    similarities = [similarities_ptr[i] for i in range(len(text))]

    return similarities

# Example usage
if __name__ == "__main__":
    result = normalize_text("Hello  World!", convert_half_space=True)
    print("Normalized:", result)