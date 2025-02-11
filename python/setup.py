import os
import subprocess
from setuptools import setup, find_packages
from setuptools.command.build_py import build_py
from wheel.bdist_wheel import bdist_wheel as _bdist_wheel

class BuildSharedLibrary(build_py):
    """Custom build command to compile Go code into a shared library."""
    def run(self):
        # Get target platform from environment variables
        goos = os.environ.get("GOOS", "linux")
        goarch = os.environ.get("GOARCH", "amd64")

        # Correct project root calculation (setup.py is in the "python" directory)
        project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
        export_go_path = os.path.join(project_root, "python", "export.go")
        output_dir = os.path.join(project_root, "python", "seperno")

        # Ensure output directory exists
        os.makedirs(output_dir, exist_ok=True)

        # Set compiler and extension based on OS
        if goos == "windows":
            cc = "x86_64-w64-mingw32-gcc"
            ext = "dll"
        elif goos == "darwin":
            cc = "clang"
            ext = "dylib"
        else:
            cc = "gcc"
            ext = "so"

        # Build command with corrected paths
        cmd = (
            f"cd {project_root} && "
            f"CC={cc} "
            f"CGO_ENABLED=1 "
            f"GOOS={goos} "
            f"GOARCH={goarch} "
            f"go build -o {output_dir}/seperno.{ext} "
            f"-buildmode=c-shared {export_go_path}"
        )

        print(f"Running Go build command: {cmd}")
        subprocess.check_call(cmd, shell=True)

        # Continue with regular build
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