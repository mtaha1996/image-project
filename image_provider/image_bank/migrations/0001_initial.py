# Generated by Django 4.2.6 on 2024-06-06 16:52

from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Access',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('engin_id', models.CharField(max_length=100, unique=True)),
                ('api_key', models.CharField(max_length=100, unique=True)),
            ],
        ),
        migrations.CreateModel(
            name='Image',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('query', models.CharField(max_length=100)),
                ('link', models.CharField(max_length=500, unique=True)),
            ],
        ),
    ]
