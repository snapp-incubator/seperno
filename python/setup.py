import os
import platform
import subprocess
from setuptools import setup, find_packages
from setuptools.command.build_py import build_py
from wheel.bdist_wheel import bdist_wheel as _bdist_wheel

class BuildSharedLibrary(build_py):
    """Custom build command to compile the Go shared library before packaging."""
    def run(self):
        system = platform.system()
        # Since setup.py is in the python/ directory, the project root is one level up.
        project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
        # export.go is inside python/ (relative to project root)
        export_go_path = os.path.join(project_root, "python", "export.go")

        # Define Go build commands for each OS:
        go_build_cmds = {
            "Linux": f"cd {project_root} && CC=gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o python/seperno/seperno.so -buildmode=c-shared {export_go_path}",
            "Darwin": f"cd {project_root} && CC=clang CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o python/seperno/seperno.dylib -buildmode=c-shared {export_go_path}",
            "Windows": f"cd {project_root} && CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o python/seperno/seperno.dll -buildmode=c-shared {export_go_path}",
        }

        if system not in go_build_cmds:
            raise RuntimeError(f"Unsupported OS: {system}")

        print("Running Go build command:")
        print(go_build_cmds[system])
        subprocess.check_call(go_build_cmds[system], shell=True)

        # Continue with the regular build process.
        build_py.run(self)

class bdist_wheel(_bdist_wheel):
    """Force the wheel to be non-pure (platform-specific)."""
    def finalize_options(self):
        _bdist_wheel.finalize_options(self)
        self.root_is_pure = False

# Choose the shared library file name based on OS.
system = platform.system()
shared_libs = {
    "Linux": "seperno.so",
    "Darwin": "seperno.dylib",
    "Windows": "seperno.dll",
}
shared_lib = shared_libs.get(system)
if not shared_lib:
    raise RuntimeError(f"Unsupported OS: {system}")

setup(
    name="seperno",
    version="1.1.5",
    author="Sepehr Sohrabpour",
    author_email="sepehrxsohrabpour@gmail.com",
    description="Python wrapper for Go-based Seperno text normalization",
    packages=find_packages(),  # Looks for packages in the same directory as setup.py (i.e. python/)
    package_data={"seperno": [shared_lib]},  # Include the shared library in the package.
    include_package_data=True,
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
    cmdclass={
        "build_py": BuildSharedLibrary,
        "bdist_wheel": bdist_wheel,  # Use our custom wheel command.
    },
)