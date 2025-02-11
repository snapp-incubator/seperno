import os
import platform
import subprocess
from setuptools import setup, find_packages
from setuptools.command.build_py import build_py

class BuildSharedLibrary(build_py):
    """ Custom build command to ensure the Go shared library is compiled. """

    def run(self):
        system = platform.system()
        project_root = os.path.abspath(os.path.dirname(__file__))  # Get root path
        python_dir = os.path.join(project_root, "python")  # Path to Python package
        export_go_path = os.path.join(python_dir, "export.go")  # Path to export.go

        # Ensure we run go build from the project root
        go_build_cmds = {
            "Linux": f"cd {project_root} && CC=gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o python/seperno/seperno.so -buildmode=c-shared {export_go_path}",
            "Darwin": f"cd {project_root} && CC=clang CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o python/seperno/seperno.dylib -buildmode=c-shared {export_go_path}",
            "Windows": f"cd {project_root} && CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o python/seperno/seperno.dll -buildmode=c-shared {export_go_path}",
        }

        if system not in go_build_cmds:
            raise RuntimeError(f"Unsupported OS: {system}")

        print(f"Running: {go_build_cmds[system]}")
        subprocess.check_call(go_build_cmds[system], shell=True)

        build_py.run(self)

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
    packages=find_packages(where="python"),
    package_dir={"": "python"},
    package_data={"seperno": [shared_lib]},  # Ensure shared library is included
    include_package_data=True,  # Include shared library in the package
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
    cmdclass={"build_py": BuildSharedLibrary},  # Run Go build before packaging
)