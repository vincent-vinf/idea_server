import os
import time

import requests
import threading
from flask import Flask, request, abort
from cairosvg import svg2png
from aip import AipContentCensor

app = Flask(__name__)
generate_url = 'https://source.boringavatars.com/beam'
img_dir = './img/'
APP_ID = '***'
API_KEY = '***'
SECRET_KEY = '***'

client = AipContentCensor(APP_ID, API_KEY, SECRET_KEY)


def baidu_content_audit(text):
    result = client.textCensorUserDefined(text)

    if result.get('error_code'):
        return False, result.get('error_msg')
    else:
        conclusionType = result.get('conclusionType')
        if conclusionType == 2:
            err_msg = []
            for i in result.get('data'):
                err_msg.append(i['msg'])
            return True, '，'.join(err_msg)
        return True, '合规'


def save_svg(path, content):
    file = open(path, 'w')
    file.write(content)
    file.close()


def generate_avatar(user_id):
    try:
        r = requests.get(generate_url)
        if r.status_code == 200:
            # mkdir directory
            if not os.path.exists(img_dir):
                os.mkdir(img_dir)
            img_png_path = '{path}avatar_{id}.png'.format(path=img_dir, id=user_id)
            img_svg_path = '{path}avatar_{id}.svg'.format(path=img_dir, id=user_id)

            thread = threading.Thread(target=save_svg,
                                      args=(img_svg_path, str(r.content, encoding='utf-8')))  # encoding remove 'b'
            thread.start()
            svg2png(bytestring=r.content, write_to=img_png_path, output_width=200,
                    output_height=200)  # the directory needs to be created manually
            return True, img_png_path
        else:
            raise Exception('Request failed!')
    except Exception as e:
        return False, e


@app.route('/audit_content', methods=['POST'])
def audit_content():
    ok, msg = baidu_content_audit(request.json['text'])
    if ok:
        return msg, 200
    else:
        abort(500, msg)


@app.route('/get_avatar_url', methods=['POST'])
def get_avatar_url():
    ok, msg = generate_avatar(request.json['userId'])
    if ok:
        return msg, 200
    else:
        abort(500, msg)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=9998)
