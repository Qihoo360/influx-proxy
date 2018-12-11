import time
import requests
import multiprocessing as mlp
from mytime import mytime


def url(precision='ns'):
    return 'http://localhost:7076/write?db=test&precision='+precision


def send(m, precision, number, init_time):
    t = init_time.t_p(precision)
    for i in range(number):
        d = m + ' fd=0 ' + str(t+i)
        r = requests.post(url(precision), d)
        if i % 500 == 0:
            time.sleep(1)
            print(m +' '+precision+' '+str(i))

def main():
    init_time = mytime(2015,1,1,0,0,0)
    thread_list = []
    thread_list.append(
        mlp.Process(target=send, args=('cpu,pr=s', 's', 100000, init_time)))
    thread_list.append(
        mlp.Process(target=send, args=('cpu,pr=h', 'h', 88888, init_time)))
    thread_list.append(
        mlp.Process(target=send, args=('mem,pr=m', 'm', 130000, init_time)))
    thread_list.append(
        mlp.Process(target=send, args=('mem,pr=ns', 'ns', 77777, init_time)))
    thread_list.append(
        mlp.Process(target=send, args=('mms', 'ms', 111111, init_time)))
    thread_list.append(
        mlp.Process(target=send, args=('mus', 'us', 122222, init_time)))
    return thread_list


if __name__ == '__main__':
    thread_list = main()
    for thr in thread_list:
        thr.start()
    for thr in thread_list:
        thr.join()
    print('write over')
