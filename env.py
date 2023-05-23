#!/usr/bin/env python3
import os

def get_env_variable(var_name):
    try:
        return os.environ[var_name]
    except KeyError:
        error_msg = f"failed to get the {var_name} environment variable"
        raise Exception(error_msg)

def read_env_file(filename):
    with open(filename) as f:
        env = {}
        for line in f:
            line = line.strip()
            if line.startswith('#') or '=' not in line:
                continue
            key, value = line.split('=', 1)
            env[key] = value
    return env

def write_env_file(filename, env):
    with open(filename, 'w') as f:
        for key, value in env.items():
            f.write(f"{key}={value}\n")

def main():
    env = read_env_file('app.env')
    env['GRPC_SERVER_ADDRESS'] = get_env_variable('GALAXY_GRPC_SERVER_ADDRESS')
    env['HTTP_SERVER_ADDRESS'] = get_env_variable('GALAXY_HTTP_SERVER_ADDRESS')
    env['TOKEN_SYMMETRIC_KEY'] = get_env_variable('GALAXY_TOKEN_SYMMETRIC_KEY')
    write_env_file('app.env', env)

if __name__ == '__main__':
    main()
