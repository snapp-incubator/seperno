import os
import platform
from setuptools import setup, find_packages

# Determine the shared library path based on the platform
system = platform.system()
if system == "Linux":
    lib_path = "seperno/linux/seperno.so"
elif system == "Darwin":  # macOS
    lib_path = "seperno/macos/seperno.dylib"
else:
    raise RuntimeError(f"Unsupported platform: {system}")

# Ensure the shared library exists
if not os.path.exists(lib_path):
    raise RuntimeError(f"Shared library not found: {lib_path}")

setup(
    name="seperno",
    version="0.0.11",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),
    package_data={
        "seperno": [lib_path],  # Include the correct shared library
    },
    include_package_data=True,
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)