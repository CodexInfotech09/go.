import os
import stat

def make_executable(file_name):
    if os.path.isfile(file_name):
        # Get the current permissions of the file
        st = os.stat(file_name)
        # Add executable permissions for the user
        os.chmod(file_name, st.st_mode | stat.S_IEXEC)
        print(f"Changed permissions to executable for: {file_name}")
    else:
        print(f"File {file_name} does not exist.")

if __name__ == "__main__":
    make_executable('bgmi')
