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


# -------- DetectPersianNumbers binding --------

# Prototype in Go:
# func DetectPersianNumbers(input *C.char, outNums **C.longlong, outStarts **C.int, outEnds **C.int, outLen *C.int)

seperno.DetectPersianNumbers.argtypes = [
    ctypes.c_char_p,                        # input string
    ctypes.POINTER(ctypes.POINTER(ctypes.c_longlong)),  # **outNums
    ctypes.POINTER(ctypes.POINTER(ctypes.c_int)),       # **outStarts
    ctypes.POINTER(ctypes.POINTER(ctypes.c_int)),       # **outEnds
    ctypes.POINTER(ctypes.c_int)            # *outLen
]
seperno.DetectPersianNumbers.restype = None   # it returns void


def detect_persian_numbers(input_text):
    """Python wrapper for the Go-based Persian number detector."""

    # Prepare output pointers
    nums_ptr = ctypes.POINTER(ctypes.c_longlong)()
    starts_ptr = ctypes.POINTER(ctypes.c_int)()
    ends_ptr = ctypes.POINTER(ctypes.c_int)()
    length = ctypes.c_int()

    # Call the CGo function
    seperno.DetectPersianNumbers(
        input_text.encode("utf-8"),
        ctypes.byref(nums_ptr),
        ctypes.byref(starts_ptr),
        ctypes.byref(ends_ptr),
        ctypes.byref(length)
    )

    n = length.value

    # Convert C arrays into Python lists
    values = [nums_ptr[i] for i in range(n)]
    start_indices = [starts_ptr[i] for i in range(n)]
    end_indices = [ends_ptr[i] for i in range(n)]

    # Important: free the memory allocated in Go (using libc free)
    libc = ctypes.CDLL("libc.so.6" if system == "Linux" else "libc.dylib")
    libc.free(nums_ptr)
    libc.free(starts_ptr)
    libc.free(ends_ptr)

    return [
        {
            "value": values[i],
            "start_index": start_indices[i],
            "end_index": end_indices[i],
        }
        for i in range(n)
    ]