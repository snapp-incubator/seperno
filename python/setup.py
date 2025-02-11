import os
import subprocess
from setuptools import setup, find_packages
from setuptools.command.build_py import build_py
from wheel.bdist_wheel import bdist_wheel as _bdist_wheel

class BuildSharedLibrary(build_py):
    """Compile Go code into a platform-specific shared library."""
    def run(self):
        # Read target platform from cibuildwheel's environment variables
        goos = os.environ.get("GOOS", "linux")
        goarch = os.environ.get("GOARCH", "amd64")
        project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
        export_go = os.path.join(project_root, "python", "export.go")
        output_dir = os.path.join(project_root, "python", "seperno")

        # Set compiler and file extension based on GOOS
        if goos == "windows":
            cc = "x86_64-w64-mingw32-gcc"
            ext = "dll"
        elif goos == "darwin":
            cc = "clang"
            ext = "dylib"
        else:
            cc = "gcc"
            ext = "so"

        # Build command
        cmd = (
            f"cd {project_root} && "
            f"CC={cc} "
            f"CGO_ENABLED=1 "
            f"GOOS={goos} "
            f"GOARCH={goarch} "
            f"go build -o {output_dir}/seperno.{ext} "
            f"-buildmode=c-shared {export_go}"
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
    package_data={"seperno": ["*.so", "*.dylib", "*.dll"]},  # Include all possible extensions
    include_package_data=True,
    cmdclass={
        "build_py": BuildSharedLibrary,
        "bdist_wheel": bdist_wheel,
    },
)