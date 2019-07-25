#!merge_and_compress_json.py

import json
import sys
import os


def get_file_names(file_dir):
    for root, dirs, files in os.walk(file_dir):
        # print(root)     # 当前目录路径
        # print(dirs)     # 当前路径下所有子目录
        # print(files)    # 当前路径下所有非目录子文件
        return files


def read_file(file_name):
    with open(file_name, "rb") as f:
        t = f.read()
        data = json.loads(t, encoding='utf8')
        return data


def merge_data(files):
    merged_data = {}
    for file in files:
        file_components = os.path.splitext(file)
        if file_components[1] == '.json':
            data = read_file(file)
            merged_data[file_components[0]] = data
    return merged_data


def write_file(data):
    print(len(data))
    with open("mergedAndCompressed.json", "w", encoding="utf8") as f:
        json.dump(data, f, ensure_ascii=False)


file_names = get_file_names("./")
final_data = merge_data(file_names)
write_file(final_data)
