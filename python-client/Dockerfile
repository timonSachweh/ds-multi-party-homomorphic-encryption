FROM python:3.11

RUN pip install gunicorn

WORKDIR /app

COPY ./requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt

COPY ./src/ /app

CMD ["sh", "-c", "gunicorn -w 2 -b 0.0.0.0:${PYTHON_PORT} 'main:app'"]