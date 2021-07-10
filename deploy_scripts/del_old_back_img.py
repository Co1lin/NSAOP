import os
image = os.popen('docker images | grep nsaop-backend').readlines()[1]
image_id = image.split()[2]
print(os.popen('docker rmi ' + image_id).readlines())
print(os.popen('docker image prune -f').readlines())