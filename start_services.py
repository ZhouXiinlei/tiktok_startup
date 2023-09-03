import os
import subprocess

def run_command(cmd, *path_parts):
    current_dir = os.getcwd()  # 获取当前的工作目录
    path = os.path.join(*path_parts)
    os.chdir(path)

    # 创建一个新的命令提示符窗口来运行Go服务
    subprocess.Popen(["cmd", "/c", "start", "cmd", "/k", cmd])

    os.chdir(current_dir)  # 返回到原始的工作目录

if __name__ == "__main__":
    # 设置Go代理
    subprocess.check_call(["go", "env", "-w", "GOPROXY=https://goproxy.cn,direct"])

    # 清理和验证go.mod
    subprocess.check_call(["go", "mod", "tidy"])

    # 启动rpc服务
    run_command("go run contact.go", "rpc", "contact")
    run_command("go run user.go", "rpc", "user")
    run_command("go run video.go", "rpc", "video")

    # 启动http服务
    run_command("go run tikstart.go", "http")

    print("All services started!")
