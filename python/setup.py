import os
import platform
from setuptools import setup, find_packages
from setuptools.command.build_ext import build_ext
from setuptools.command.install import install
import subprocess

class BuildSharedLibrary(build_ext):
    """ Custom build command to ensure the Go shared library is compiled. """

    def run(self):
        system = platform.system()
        go_build_cmds = {
            "Linux": "CC=gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o seperno/seperno.so -buildmode=c-shared export.go",
            "Darwin": "CC=clang CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o seperno/seperno.dylib -buildmode=c-shared export.go",
            "Windows": "CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o seperno/seperno.dll -buildmode=c-shared export.go",
        }

        if system not in go_build_cmds:
            raise RuntimeError(f"Unsupported OS: {system}")

        print(f"Running: {go_build_cmds[system]}")
        subprocess.check_call(go_build_cmds[system], shell=True)

        build_ext.run(self)

# Detect OS and determine the correct shared library file
system = platform.system()
shared_libs = {
    "Linux": "seperno/seperno.so",
    "Darwin": "seperno/seperno.dylib",
    "Windows": "seperno/seperno.dll",
}
shared_lib = shared_libs.get(system, "")

if not shared_lib:
    raise RuntimeError(f"Unsupported OS: {system}")

setup(
    name="seperno",
    version="1.1.5",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),
    package_data={"seperno": [shared_lib]},  # Ensure shared library is included
    include_package_data=True,  # Include shared library in the package
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
    cmdclass={"build_ext": BuildSharedLibrary},  # Run Go build before packaging
)