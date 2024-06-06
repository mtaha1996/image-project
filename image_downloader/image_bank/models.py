from django.db import models
from django.conf import settings


class Image(models.Model):
    query = models.CharField(max_length=100)
    link = models.CharField(max_length=500, unique=True)


class Access(models.Model):
    engin_id = models.CharField(max_length=100, unique=True)
    api_key = models.CharField(max_length=100, unique=True)
