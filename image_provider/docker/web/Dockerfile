# Use the official Python image as the base image
FROM python

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Set work directory
WORKDIR /code

# Install dependencies
COPY ./image_provider/requirements.txt /code/
RUN pip install --no-cache-dir -r requirements.txt

# Copy project
COPY ./image_provider /code/

CMD ["python", "manage.py runserver 0.0.0.0:8000"]
