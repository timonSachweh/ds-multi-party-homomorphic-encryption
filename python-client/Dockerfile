FROM pytorch/pytorch:2.7.1-cuda12.8-cudnn9-runtime
WORKDIR /app

#COPY ./requirements.txt /app/requirements.txt
#RUN pip install -r requirements.txt
RUN python3 -m pip install numpy flask

COPY ./src/ /app

CMD ["python3", "main.py"]
#CMD ["sh", "-c", "gunicorn --preload -w 1 -b 0.0.0.0:${PYTHON_PORT} 'main:app'"]
