FROM python:3.11

RUN pip install gunicorn

WORKDIR /app

COPY ./requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt

COPY ./src/ /app

CMD ["python3", "main.py"]
#CMD ["sh", "-c", "gunicorn --preload -w 1 -b 0.0.0.0:${PYTHON_PORT} 'main:app'"]