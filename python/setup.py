import os
import platform
from setuptools import setup, find_packages

# Detect OS to include the right shared library
system = platform.system()
if system == "Linux":
    shared_lib = "seperno.so"
elif system == "Darwin":  # macOS
    shared_lib = "seperno.dylib"
elif system == "Windows":
    shared_lib = "seperno.dll"
else:
    raise RuntimeError(f"Unsupported OS: {system}")

setup(
    name="seperno",
    version="1.1.5",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),
    package_data={"seperno": [shared_lib]},
    include_package_data=True,
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)