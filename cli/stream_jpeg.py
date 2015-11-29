import io
import socket
import struct
import time
import picamera
import sys

client_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
client_socket.setblocking(0)

try:
    with picamera.PiCamera() as camera:
        camera.resolution = (640, 480)
        camera.framerate = 10 
        time.sleep(2)
        start = time.time()
        stream = io.BytesIO()
        # Use the video-port for captures...
        for foo in camera.capture_continuous(stream, 'jpeg',
                                             use_video_port=True, quality=10):
            #connection.write(struct.pack('<L', stream.tell()))
            #connection.flush()
            size = stream.tell()

            stream.seek(0)

            try:
                client_socket.sendto(stream.read(), ('__YOUR SERVER__', 3000))
            except socket.error as e:
                print "IOError: %s" % e
            except:
                print "Unexpected error:", sys.exc_info()[0]
            stream.seek(0)
            stream.truncate()

finally:
    print "Nu lukker vi..."
    #connection.close()
    #client_socket.close()