import os
import subprocess
from setuptools import setup, find_packages
from setuptools.command.build_py import build_py
from wheel.bdist_wheel import bdist_wheel as _bdist_wheel
import sys
import os

goos_arg = None
goarch_arg = None
cc_arg = None
new_argv = []

for arg in sys.argv:
    if arg.startswith("--goos="):
        goos_arg = arg.split("=")[1]
    elif arg.startswith("--goarch="):
        goarch_arg = arg.split("=")[1]
    elif arg.startswith("--cc="):
        cc_arg = arg.split("=")[1]
    else:
        new_argv.append(arg)
sys.argv = new_argv

if goos_arg:
    os.environ["GOOS"] = goos_arg
if goarch_arg:
    os.environ["GOARCH"] = goarch_arg
if cc_arg:
    os.environ["CC"] = cc_arg

class BuildSharedLibrary(build_py):
    def run(self):
        goos = os.environ.get("GOOS", "linux")
        goarch = os.environ.get("GOARCH", "amd64")

        # CORRECTED: Path calculation (setup.py is in "python" dir)
        project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
        export_go_path = os.path.join(project_root, "python", "export.go")
        output_dir = os.path.join(project_root, "python", "seperno")

        # Create output directory if needed
        os.makedirs(output_dir, exist_ok=True)

        # Set compiler and extension
        if goos == "windows":
            cc = "x86_64-w64-mingw32-gcc"
            ext = "dll"
        elif goos == "darwin":
            cc = "clang"
            ext = "dylib"
        else:
            cc = "gcc"
            ext = "so"

        # CORRECTED: Build command (uses project_root, not "/")
        cmd = (
            f"cd {project_root} && "  # <-- Critical fix here
            f"CC={cc} "
            f"CGO_ENABLED=1 "
            f"GOOS={goos} "
            f"GOARCH={goarch} "
            f"go build -o {output_dir}/seperno.{ext} "
            f"-buildmode=c-shared {export_go_path}"
        )

        print(f"Running Go build command: {cmd}")
        subprocess.check_call(cmd, shell=True)
        super().run()

class bdist_wheel(_bdist_wheel):
    """Mark wheel as platform-specific."""
    def finalize_options(self):
        super().finalize_options()
        self.root_is_pure = False

setup(
    name="seperno",
    version="1.1.5",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),
    package_data={"seperno": ["*.so", "*.dylib", "*.dll"]},
    include_package_data=True,
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
    cmdclass={
        "build_py": BuildSharedLibrary,
        "bdist_wheel": bdist_wheel,
    },
)