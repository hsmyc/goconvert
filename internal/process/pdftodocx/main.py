import os
import argparse
from pdf2docx import parse


def convert_pdf_to_docx(pdf_file, docx_file):
    parse(pdf_file, docx_file)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Convert PDF to DOCX")
    parser.add_argument("pdf_file", help="Path to the input PDF file")
    parser.add_argument("docx_file", help="Path to the output DOCX file")

    args = parser.parse_args()

    convert_pdf_to_docx(args.pdf_file, args.docx_file)
