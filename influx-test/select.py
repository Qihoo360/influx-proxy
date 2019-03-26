import requests
import time
from mytime import mytime


def time_TZ(str_t, precision):
    """

    :type str_t: str
    """
    if '.' in str_t:
        tmp = str_t.split('.')[0]
    else:
        tmp = str_t[:-1]
    now = time.strptime(tmp, '%Y-%m-%dT%X')
    nowt = int(time.mktime(now))
    if precision == 'h':
        nowt = nowt / 3600
    elif precision == 'm':
        nowt = nowt / 60
    elif precision == 'ms':
        nowt = nowt * 1000
        w = 3
    elif precision == 'us':
        nowt = nowt * 1000000
        w = 6
    elif precision == 'ns':
        nowt = nowt * 1000000000
        w = 9

    if '.' in str_t:
        tmp2 = str_t.split('.')[1][:-1]
        nowt = nowt + int(tmp2) * (10 ** (w - len(tmp2)))

    return nowt




def check(m, precision, number, init_time):
    url = 'http://localhost:7076/query'
    payload = {}
    payload['db'] = 'test'
    tmp = 'select fd from ' + m.split(',')[0]
    if len(m.split(','))>1:
        tmp2 = m.split(',')[1]
        pr=tmp2.split('=')[1]
        tmp=tmp + " where pr = '%s'" % pr
    payload['q']=tmp

    r = requests.get(url, params=payload)
    if r.status_code != 200:
        print('http error! expect: 200 actual: %s' % r.status_code)
        return
    body = r.json()['results'][0]['series'][0]['values']
    if len(body) != number:
        print('number error! expect: %s actual: %s' % (number, len(body)))
        return
    tl = [x[0] for x in body]

    for i in range(0, number):
        if init_time.t_p(precision) + i != time_TZ(tl[i], precision):
            print('time error! expect: %s actual: %s' %
                  (init_time.t_p(precision) + i,
                   time_TZ(tl[i], precision)))
            return
        if i % 5000 == 0:
            print('check %s' % i)
    return

def main():
    init_time = mytime(2014, 12, 31, 16, 0, 0 )
    check('cpu,pr=s', 's', 100000, init_time)
    check('cpu,pr=h', 'h', 88888, init_time)
    check('mem,pr=m', 'm', 130000, init_time)
    check('mem,pr=ns', 'ns', 77777, init_time)
    check('mms', 'ms', 111111, init_time)
    check('mus', 'us', 122222, init_time)


if __name__ == '__main__':
    main()