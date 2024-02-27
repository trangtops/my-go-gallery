import os
import json

item_per_page = 50
db_name = "db"
db = {}

# Direc = input(r"Enter the path of the folder: ")
Direc = "D:\Photo\p"
print(f"folders in the directory: {Direc}")

folders = os.listdir(Direc)
folders = [f for f in folders if not os.path.isfile(Direc+'/'+f)] #Filtering only the folders.

gallery = {}
for i in range(len(folders)):
    alblum = {}
    alblum["name"] = folders[i]
    path = os.path.join(Direc, folders[i])
    alblum["path"] = path
    files = os.listdir(path)
    alblum["thumbnail"] = files[0]
    gallery[str(i)] = alblum

gallery[str(len(folders))] = {
    "name": "twitter",
    "path": "D:\\Photo\\twitter_media_harvest",
    "thumbnail": os.listdir("D:/Photo/twitter_media_harvest")[0]
}
gallery[str(len(folders) + 1)] = {
    "name": "twitter",
    "path": "D:\\Photo\\t2",
    "thumbnail": os.listdir("D:/Photo/t2")[0]
}
# alblum["pages"] = {}

# file_len = len(folders)
# for i in range(int(file_len/item_per_page)):
#     page = []
#     for j in range(item_per_page):
#         file_index = (i*item_per_page)+j
#         if file_index > file_len-1:
#             break 
#         page.append(folders[file_index])

#     alblum["pages"][i] = page
# db[Direc] = alblum

with open(db_name, 'w', encoding='utf8') as f:
    json.dump(gallery, f, ensure_ascii=False)

#os.getcwd() gives us the current working directory, and os.listdir lists the director
