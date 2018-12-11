import time


class mytime:
    def fz(self, x):
        # front zero
        if x / 10 >= 1:
            return str(x)
        else:
            return '0'+str(x)

    def __init__(self, Y, M, D, h, m, s, ms=0, us=0, ns=0):
        self.format_time = str(Y) + '-' + self.fz(M) + '-' + self.fz(D) + ' ' + \
            self.fz(h) + ':' + self.fz(m) + ':' + self.fz(s)
        self.format = '%Y-%m-%d %X'
        self.struct_time = time.strptime(self.format_time, self.format)
        self.timestamp = int(time.mktime(self.struct_time))
        self.ms = ms
        self.us = us
        self.ns = ns

    def t_h(self):
        return self.timestamp / 3600

    def t_m(self):
        return self.timestamp / 60

    def t_s(self):
        return self.timestamp

    def t_ms(self):
        return self.timestamp * 1000 + self.ms

    def t_us(self):
        return self.timestamp * 1000000 + self.us

    def t_ns(self):
        return self.timestamp * 1000000000 + self.ns

    def after(self, sec):
        # offer a fake time
        # just ensure its timestamp to be correct
        a = mytime(2000,1,1,1,1,1, self.ms, self.us, self.ns)
        a.timestamp = self.timestamp + sec
        return a

    def t_p(self, precision):
        td = {
            'h': self.t_h(),
            'm': self.t_m(),
            's': self.t_s(),
            'ms': self.t_ms(),
            'us': self.t_us(),
            'ns': self.t_ns(),
        }
        return td[precision]