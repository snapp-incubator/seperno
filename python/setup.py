import platform
from setuptools import setup, find_packages

# Determine the shared library extension based on the platform
system = platform.system()
if system == "Linux":
    lib_extension = ".so"
elif system == "Darwin":  # macOS
    lib_extension = ".dylib"
else:
    raise RuntimeError(f"Unsupported platform: {system}")

setup(
    name="seperno",
    version="1.1.5",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),
    package_data={
        "seperno": [f"*{lib_extension}"],
    },
    include_package_data=True,
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)