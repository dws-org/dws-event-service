# Wir starten mit einem Go‑Basisimage
FROM golang:1.23.1

# Benötigt, damit "make" auf dem Image vorhanden ist
RUN apt-get update && apt-get install -y make

# Arbeitsverzeichnis im Container
WORKDIR /app

# Den gesamten Quellcode in das Arbeitsverzeichnis kopieren
COPY . /app

# Prisma-Client generieren (falls in Ihrem Projekt benötigt)
RUN ["make", "generate"]

# Port freigeben – wählen Sie den Port, auf dem Ihr Server läuft
EXPOSE 6906

# Standardbefehl: Makefile-Task "run" ausführen
CMD ["make", "run"]
