import random
import string


class HelperClass(object):

    @staticmethod
    def random_word_special_signs_included(word_length):
        chars = string.ascii_letters + string.digits + '!@#$%^&*()/\{}[]<>'
        return ''.join(random.sample(chars, word_length))

    @staticmethod
    def random_word_letters_only(word_length):
        chars = string.ascii_letters
        return ''.join(random.sample(chars, word_length))

    @staticmethod
    def create_random_email():
        first_part = random_word(15)
        second_part = random_word(5)
        return f'{first_part}@{second_part}.com'
