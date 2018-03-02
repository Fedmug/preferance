#!/usr/bin/env python2
# -*- coding: utf-8 -*-

from random import shuffle, sample, choice, randint
import numpy as np
from copy import copy
import argparse

CARDS_SIZE = 32
SUIT_SIZE = 8
SUITS_NUMBER = 4
SUIT_SYMBOLS = (u"\u2660", u"\u2663", u"\u2666", u"\u2665", u"БК", u"МИЗЕР")
CARD_RANKS = ('7', '8', '9', '10', 'J', 'Q', 'K', 'A')
SUIT_NAMES = {'spades': 0, 'clubs': 1, 'diamonds': 2, 'hearts': 3, 'nt': 4}

INITIAL_TRICKS_BOUND = 4
INITIAL_MISERE_BOUND = 1
TRICKS_STEP = 0.5
MISERE_STEP = 0.5


NON_TRUMP_WITH_LONG_TRUMP = 0
NON_TRUMP_WITH_SHORT_TRUMP = 1
TRUMP = 2
MISERE = 3
TRICKS_ESTIMATIONS = (

    # [tricks if suit is not trump but there is long enough trump suit in hand,
    #  tricks if suit is not trump and there is no long enough trump suit in hand,
    #  tricks if suit is trump,
    #  misere coefficient]

    # empty suit
    (0, 0, 0, 0),
    # 0x01 = '7'
    (0, 0, 0, 0),
    # 0x02 = '8'
    (0, 0, 0, 0.1),
    # 0x03 = '87'
    (0, 0, 0, 0),
    # 0x04 = '9'
    (0, 0, 0, 1),
    # 0x05 = '97'
    (0, 0, 0, 0),
    # 0x06 = '98'
    (0, 0, 0, 1),
    # 0x07 = '987'
    (0, 0, 0.1, 0),
    # 0x08 = '10'
    (0, 0, 0, 1.5),
    # 0x09 = '107'
    (0, 0, 0, 0.2),
    # 0x0A = '108'
    (0, 0, 0, 1.1),
    # 0x0B = '1087'
    (0, 0, 0.1, 0),
    # 0x0C = '109'
    (0, 0, 0, 2),
    # 0x0D = '1097'
    (0, 0, 0.1, 0),
    # 0x0E = '1098'
    (0, 0, 0.1, 1.6),
    # 0x0F = '10987'
    (0.8, 0, 1.35, 0),
    # 0x10 = 'J'
    (0, 0, 0, 1.8),
    # 0x11 = 'J7'
    (0, 0, 0, 0.5),
    # 0x12 = 'J8'
    (0, 0, 0, 1.3),
    # 0x13 = 'J87'
    (0, 0, 0.2, 0),
    # 0x14 = 'J9'
    (0, 0, 0, 2.2),
    # 0x15 = 'J97'
    (0, 0, 0.2, 0),
    # 0x16 = 'J98'
    (0, 0, 0.2, 1.7),
    # 0x17 = 'J987'
    (0.8, 0, 1.4, 0),
    # 0x18 = 'J10'
    (0, 0, 0, 2.7),
    # 0x19 = 'J107'
    (0, 0, 0.25, 0.4),
    # 0x1A = 'J108'
    (0, 0, 0.25, 1.75),
    # 0x1B = 'J1087'
    (0.8, 0, 1.4, 0),
    # 0x1C = 'J109'
    (0, 0, 0.3, 3.1),
    # 0x1D = 'J1097'
    (0.8, 0, 1.4, 0),
    # 0x1E = 'J1098'
    (0.82, 0, 1.41, 2.5),
    # 0x1F = 'J10987'
    (2.5, 0, 2.8, 0),
    # 0x20 = 'Q'
    (0, 0, 0, 1.9),
    # 0x21 = 'Q7'
    (0, 0, 0.1, 0.7),
    # 0x22 = 'Q8'
    (0, 0, 0.1, 1.5),
    # 0x23 = 'Q87'
    (0.01, 0.01, 0.3, 0.2),
    # 0x24 = 'Q9'
    (0, 0, 0.1, 2.3),
    # 0x25 = 'Q97'
    (0.01, 0.01, 0.3, 0.2),
    # 0x26 = 'Q98'
    (0.01, 0.01, 0.3, 1.8),
    # 0x27 = 'Q987'
    (0.8, 0, 1.66, 0),
    # 0x28 = 'Q10'
    (0, 0, 0.1, 2.9),
    # 0x29 = 'Q107'
    (0.01, 0.01, 0.35, 0.6),
    # 0x2A = 'Q108'
    (0.01, 0.01, 0.35, 1.85),
    # 0x2B = 'Q1087'
    (0.8, 0, 1.67, 0),
    # 0x2C = 'Q109'
    (0.05, 0.05, 0.5, 3.2),
    # 0x2D = 'Q1097'
    (0.82, 0, 1.67, 0),
    # 0x2E = 'Q1098'
    (0.85, 0.01, 1.7, 2.6),
    # 0x2F = 'Q10987'
    (2.5, 0, 2.9, 0),
    # 0x30 = 'QJ'
    (0, 0, 0.1, 3),
    # 0x31 = 'QJ7'
    (0.08, 0.08, 0.9, 1.3),
    # 0x32 = 'QJ8'
    (0.09, 0.09, 0.92, 2.4),
    # 0x33 = 'QJ87'
    (0.9, 0, 1.72, 0),
    # 0x34 = 'QJ9'
    (0.2, 0.2, 0.98, 3.3),
    # 0x35 = 'QJ97'
    (0.91, 0, 1.72, 0),
    # 0x36 = 'QJ98'
    (0.92, 0, 1.74, 2.65),
    # 0x37 = 'QJ987'
    (2.5, 0, 2.98, 0),
    # 0x38 = 'QJ10'
    (0.3, 0.3, 1, 3.5),
    # 0x39 = 'QJ107'
    (1.25, 0.05, 1.75, 0.65),
    # 0x3A = 'QJ108'
    (1.3, 0.05, 1.75, 2.7),
    # 0x3B = 'QJ1087'
    (2.5, 0, 3.03, 0),
    # 0x3C = 'QJ109'
    (1.35, 0.05, 2.02, 3.9),
    # 0x3D = 'QJ1097'
    (2.5, 0, 3.03, 0),
    # 0x3E = 'QJ1098'
    (2.5, 0, 3.03, 3.4),
    # 0x3F = 'QJ10987'
    (4.5, 0, 4.5, 0),
    # 0x40 = 'K'
    (0.01, 0.01, 0.01, 2),
    # 0x41 = 'K7'
    (0.25, 0.25, 0.5, 0.8),
    # 0x42 = 'K8'
    (0.25, 0.25, 0.5, 1.6),
    # 0x43 = 'K87'
    (0.33, 0.33, 1, 0.35),
    # 0x44 = 'K9'
    (0.25, 0.25, 0.5, 2.5),
    # 0x45 = 'K97'
    (0.33, 0.33, 1, 0.35),
    # 0x46 = 'K98'
    (0.33, 0.33, 1, 2),
    # 0x47 = 'K987'
    (0.9, 0, 1.7, 0),
    # 0x48 = 'K10'
    (0.25, 0.25, 0.5, 3),
    # 0x49 = 'K107'
    (0.35, 0.35, 1.1, 1.1),
    # 0x4A = 'K108'
    (0.35, 0.35, 1.1, 2.1),
    # 0x4B = 'K1087'
    (0.9, 0, 1.71, 0),
    # 0x4C = 'K109'
    (0.4, 0.4, 1.2, 3.4),
    # 0x4D = 'K1097'
    (0.91, 0.01, 1.73, 0),
    # 0x4E = 'K1098'
    (0.92, 0.01, 1.75, 2.8),
    # 0x4F = 'K10987'
    (2.7, 0, 3.2, 0),
    # 0x50 = 'KJ'
    (0.36, 0.36, 0.8, 3.2),
    # 0x51 = 'KJ7'
    (0.52, 0.52, 1.25, 1.4),
    # 0x52 = 'KJ8'
    (0.52, 0.52, 1.3, 2.5),
    # 0x53 = 'KJ87'
    (1, 0.02, 1.74, 0),
    # 0x54 = 'KJ9'
    (0.6, 0.6, 1.65, 3.6),
    # 0x55 = 'KJ97'
    (1, 0.03, 1.76, 0),
    # 0x56 = 'KJ98'
    (1.02, 0.03, 1.78, 2.85),
    # 0x57 = 'KJ987'
    (3, 0, 3.5, 0),
    # 0x58 = 'KJ10'
    (0.8, 0.8, 1.8, 4),
    # 0x59 = 'KJ107'
    (1.5, 0.03, 2.2, 0.8),
    # 0x5A = 'KJ108'
    (1.55, 0.03, 2.3, 2.95),
    # 0x5B = 'KJ1087'
    (3.1, 0, 3.65, 0),
    # 0x5C = 'KJ109'
    (1.6, 0.03, 2.4, 4.3),
    # 0x5D = 'KJ1097'
    (3.1, 0, 3.65, 0),
    # 0x5E = 'KJ1098'
    (3.1, 0, 3.65, 3.5),
    # 0x5F = 'KJ10987'
    (4.5, 0, 4.5, 0),
    # 0x60 = 'KQ'
    (0.9, 0.9, 1, 3.8),
    # 0x61 = 'KQ7'
    (0.95, 0.95, 1.3, 1.7),
    # 0x62 = 'KQ8'
    (0.95, 0.95, 1.3, 3.1),
    # 0x63 = 'KQ87'
    (2, 0.05, 2.6, 0.25),
    # 0x64 = 'KQ9'
    (1, 1, 1.5, 3.85),
    # 0x65 = 'KQ97'
    (2, 0.05, 2.6, 0.25),
    # 0x66 = 'KQ98'
    (2, 0.05, 2.6, 2.9),
    # 0x67 = 'KQ987'
    (3.5, 0, 3.82, 0),
    # 0x68 = 'KQ10'
    (1.1, 1.1, 1.85, 4.2),
    # 0x69 = 'KQ107'
    (2.05, 0.05, 2.65, 0.9),
    # 0x6A = 'KQ108'
    (2.1, 0.05, 2.7, 3),
    # 0x6B = 'KQ1087'
    (3.5, 0, 3.82, 0),
    # 0x6C = 'KQ109'
    (2.15, 0.05, 2.8, 4.5),
    # 0x6D = 'KQ1097'
    (3.5, 0, 3.82, 0),
    # 0x6E = 'KQ1098'
    (3.5, 0, 3.82, 3.65),
    # 0x6F = 'KQ10987'
    (5, 0, 5, 0),
    # 0x70 = 'KQJ'
    (1.5, 1.5, 2, 4.4),
    # 0x71 = 'KQJ7'
    (2.1, 0.07, 2.98, 1.65),
    # 0x72 = 'KQJ8'
    (2.1, 0.07, 2.99, 3.2),
    # 0x73 = 'KQJ87'
    (3.2, 0, 4, 0),
    # 0x74 = 'KQJ9'
    (2.1, 0.07, 3, 4.4),
    # 0x75 = 'KQJ97'
    (3.2, 0, 4, 0),
    # 0x76 = 'KQJ98'
    (3.2, 0, 4, 3.65),
    # 0x77 = 'KQJ987'
    (4.5, 0, 5, 0),
    # 0x78 = 'KQJ10'
    (2.1, 0.07, 3, 4.6),
    # 0x79 = 'KQJ107'
    (3.2, 0, 4, 0.95),
    # 0x7A = 'KQJ108'
    (3.2, 0, 4, 4.1),
    # 0x7B = 'KQJ1087'
    (4.5, 0, 5, 0),
    # 0x7C = 'KQJ109'
    (3.2, 0, 4, 5),
    # 0x7D = 'KQJ1097'
    (4.5, 0, 5, 0),
    # 0x7E = 'KQJ1098'
    (4.5, 0, 5, 4.5),
    # 0x7F = 'KQJ10987'
    (6, 0, 6, 0),
    # 0x80 = 'A'
    (1, 1, 1, 2),
    # 0x81 = 'A7'
    (0.99, 0.99, 1.01, 0.9),
    # 0x82 = 'A8'
    (0.99, 0.99, 1.01, 2.1),
    # 0x83 = 'A87'
    (0.98, 0.98, 1.2, 0.45),
    # 0x84 = 'A9'
    (0.99, 0.99, 1.01, 2.9),
    # 0x85 = 'A97'
    (0.98, 0.98, 1.2, 0.45),
    # 0x86 = 'A98'
    (0.98, 0.98, 1.2, 2.75),
    # 0x87 = 'A987'
    (2.1, 0.95, 2.5, 0.087),
    # 0x88 = 'A10'
    (0.99, 0.99, 1.01, 3.3),
    # 0x89 = 'A107'
    (0.98, 0.98, 1.2, 1.1),
    # 0x8A = 'A108'
    (0.98, 0.98, 1.2, 2.85),
    # 0x8B = 'A1087'
    (2.1, 0.95, 2.5, 0.087),
    # 0x8C = 'A109'
    (0.98, 0.98, 1.25, 3.5),
    # 0x8D = 'A1097'
    (2.1, 0.95, 2.5, 0.087),
    # 0x8E = 'A1098'
    (2.1, 0.95, 2.5, 3.1),
    # 0x8F = 'A10987'
    (3.79, 0.7, 3.95, 0),
    # 0x90 = 'AJ'
    (0.99, 0.99, 1.02, 3.4),
    # 0x91 = 'AJ7'
    (1.02, 1.02, 1.3, 2.2),
    # 0x92 = 'AJ8'
    (1.02, 1.02, 1.3, 3.2),
    # 0x93 = 'AJ87'
    (2.1, 0.95, 2.5, 0.087),
    # 0x94 = 'AJ9'
    (1.05, 1.05, 1.35, 3.7),
    # 0x95 = 'AJ97'
    (2.1, 0.95, 2.5, 0.087),
    # 0x96 = 'AJ98'
    (2.1, 0.98, 2.5, 3.1),
    # 0x97 = 'AJ987'
    (3.79, 0.75, 4, 0),
    # 0x98 = 'AJ10'
    (1.2, 1.2, 1.5, 3.9),
    # 0x99 = 'AJ107'
    (2.15, 0.95, 2.6, 0.95),
    # 0x9A = 'AJ108'
    (2.15, 0.95, 2.6, 3.2),
    # 0x9B = 'AJ1087'
    (3.79, 0.75, 4, 0),
    # 0x9C = 'AJ109'
    (2.18, 0.95, 2.66, 4),
    # 0x9D = 'AJ1097'
    (3.79, 0.75, 4, 0),
    # 0x9E = 'AJ1098'
    (3.79, 0.75, 4, 3.5),
    # 0x9F = 'AJ10987'
    (5.48, 0.4, 5.55, 0),
    # 0xA0 = 'AQ'
    (1.25, 1.25, 1.5, 3.5),
    # 0xA1 = 'AQ7'
    (1.4, 1.4, 2, 2.5),
    # 0xA2 = 'AQ8'
    (1.4, 1.4, 2, 3.3),
    # 0xA3 = 'AQ87'
    (2.3, 0.98, 2.95, 0.27),
    # 0xA4 = 'AQ9'
    (1.5, 1.5, 2.1, 3.9),
    # 0xA5 = 'AQ97'
    (2.3, 0.98, 2.95, 0.27),
    # 0xA6 = 'AQ98'
    (2.3, 0.98, 2.98, 3.4),
    # 0xA7 = 'AQ987'
    (3.9, 0.75, 4.15, 0),
    # 0xA8 = 'AQ10'
    (1.7, 1.7, 2.3, 4.2),
    # 0xA9 = 'AQ107'
    (2.3, 1, 3.1, 1),
    # 0xAA = 'AQ108'
    (2.3, 1, 3.1, 3.5),
    # 0xAB = 'AQ1087'
    (3.9, 0.75, 4.15, 0),
    # 0xAC = 'AQ109'
    (3.9, 1, 4.15, 4.7),
    # 0xAD = 'AQ1097'
    (3.9, 0.75, 4.15, 0),
    # 0xAE = 'AQ1098'
    (3.9, 0.75, 4.15, 3.7),
    # 0xAF = 'AQ10987'
    (5.48, 0.4, 5.55, 0),
    # 0xB0 = 'AQJ'
    (1.9, 1.3, 2.6, 4.3),
    # 0xB1 = 'AQJ7'
    (3.1, 1.1, 3.3, 2.3),
    # 0xB2 = 'AQJ8'
    (3.1, 1.1, 3.3, 4),
    # 0xB3 = 'AQJ87'
    (4.1, 0.75, 4.4, 0),
    # 0xB4 = 'AQJ9'
    (3.15, 1.1, 3.3, 4.8),
    # 0xB5 = 'AQJ97'
    (4.1, 0.75, 4.4, 0),
    # 0xB6 = 'AQJ98'
    (4.1, 0.75, 4.4, 3.7),
    # 0xB7 = 'AQJ987'
    (5.48, 0.4, 5.55, 0),
    # 0xB8 = 'AQJ10'
    (3.2, 1.1, 3.3, 5),
    # 0xB9 = 'AQJ107'
    (4.1, 0.75, 4.4, 0.95),
    # 0xBA = 'AQJ108'
    (4.1, 0.75, 4.4, 3.8),
    # 0xBB = 'AQJ1087'
    (5.48, 0.4, 5.55, 0),
    # 0xBC = 'AQJ109'
    (4.1, 0.75, 4.4, 5.5),
    # 0xBD = 'AQJ1097'
    (5.48, 0.4, 5.55, 0),
    # 0xBE = 'AQJ1098'
    (5.48, 0.4, 5.55, 4.6),
    # 0xBF = 'AQJ10987'
    (7, 0, 7, 0),
    # 0xC0 = 'AK'
    (1.9, 1.9, 2, 3.7),
    # 0xC1 = 'AK7'
    (1.8, 1.8, 2.3, 2.7),
    # 0xC2 = 'AK8'
    (1.8, 1.8, 2.3, 3.6),
    # 0xC3 = 'AK87'
    (3.3, 1.2, 3.5, 0.52),
    # 0xC4 = 'AK9'
    (1.8, 1.8, 2.3, 4.2),
    # 0xC5 = 'AK97'
    (3.3, 1.2, 3.5, 0.52),
    # 0xC6 = 'AK98'
    (3.3, 1.2, 3.5, 3.5),
    # 0xC7 = 'AK987'
    (4.79, 0.75, 4.86, 0),
    # 0xC8 = 'AK10'
    (1.85, 1.8, 2.4, 4.5),
    # 0xC9 = 'AK107'
    (3.3, 1.2, 3.5, 1.2),
    # 0xCA = 'AK108'
    (3.3, 1.2, 3.5, 3.8),
    # 0xCB = 'AK1087'
    (4.79, 0.75, 4.86, 0),
    # 0xCC = 'AK109'
    (3.3, 1.2, 3.5, 4.6),
    # 0xCD = 'AK1097'
    (4.79, 0.75, 4.86, 0),
    # 0xCE = 'AK1098'
    (4.79, 0.75, 4.86, 3.9),
    # 0xCF = 'AK10987'
    (6, 0.4, 6, 0),
    # 0xD0 = 'AKJ'
    (1.95, 1.95, 2.9, 4.7),
    # 0xD1 = 'AKJ7'
    (3.4, 1.25, 3.6, 2.4),
    # 0xD2 = 'AKJ8'
    (3.4, 1.25, 3.6, 4),
    # 0xD3 = 'AKJ87'
    (4.79, 0.75, 4.86, 0),
    # 0xD4 = 'AKJ9'
    (3.4, 1.3, 3.62, 4.7),
    # 0xD5 = 'AKJ97'
    (4.79, 0.75, 4.86, 0),
    # 0xD6 = 'AKJ98'
    (4.1, 0.75, 4.4, 3.9),
    # 0xD7 = 'AKJ987'
    (6, 0.4, 6, 0),
    # 0xD8 = 'AKJ10'
    (3.42, 1.3, 3.64, 5),
    # 0xD9 = 'AKJ107'
    (4.79, 0.75, 4.86, 1),
    # 0xDA = 'AKJ108'
    (4.79, 0.75, 4.86, 4),
    # 0xDB = 'AKJ1087'
    (6, 0.4, 6, 0),
    # 0xDC = 'AKJ109'
    (4.79, 0.75, 4.86, 5.1),
    # 0xDD = 'AKJ1097'
    (6, 0.4, 6, 0),
    # 0xDE = 'AKJ1098'
    (6, 0.4, 6, 4.6),
    # 0xDF = 'AKJ10987'
    (7, 0, 7, 0),
    # 0xE0 = 'AKQ'
    (2.7, 2.7, 3, 4),
    # 0xE1 = 'AKQ7'
    (3.7, 1.8, 3.913, 2.8),
    # 0xE2 = 'AKQ8'
    (3.7, 1.8, 3.913, 3.9),
    # 0xE3 = 'AKQ87'
    (5, 0.75, 5, 0.2),
    # 0xE4 = 'AKQ9'
    (3.7, 1.8, 3.913, 4.6),
    # 0xE5 = 'AKQ97'
    (5, 0.75, 5, 0.2),
    # 0xE6 = 'AKQ98'
    (5, 0.75, 5, 4),
    # 0xE7 = 'AKQ987'
    (6, 0.4, 6, 0),
    # 0xE8 = 'AKQ10'
    (3.72, 1.8, 3.913, 4.8),
    # 0xE9 = 'AKQ107'
    (5, 0.75, 5, 1),
    # 0xEA = 'AKQ108'
    (5, 0.75, 5, 4.2),
    # 0xEB = 'AKQ1087'
    (6, 0.4, 6, 0),
    # 0xEC = 'AKQ109'
    (5, 0.75, 5, 5.2),
    # 0xED = 'AKQ1097'
    (6, 0.4, 6, 0),
    # 0xEE = 'AKQ1098'
    (6, 0.4, 6, 4.8),
    # 0xEF = 'AKQ10987'
    (7, 0, 7, 0),
    # 0xF0 = 'AKQJ'
    (3.8, 1.8, 4, 5),
    # 0xF1 = 'AKQJ7'
    (5, 0.75, 5, 2.2),
    # 0xF2 = 'AKQJ8'
    (5, 0.75, 5, 4.8),
    # 0xF3 = 'AKQJ87'
    (6, 0.4, 6, 0),
    # 0xF4 = 'AKQJ9'
    (5, 0.75, 5, 5.3),
    # 0xF5 = 'AKQJ97'
    (6, 0.4, 6, 0),
    # 0xF6 = 'AKQJ98'
    (6, 0.4, 6, 4.8),
    # 0xF7 = 'AKQJ987'
    (7, 0, 7, 0),
    # 0xF8 = 'AKQJ10'
    (5, 0.75, 5, 5.5),
    # 0xF9 = 'AKQJ107'
    (6, 0.4, 6, 1),
    # 0xFA = 'AKQJ108'
    (6, 0.4, 6, 4.8),
    # 0xFB = 'AKQJ1087'
    (7, 0, 7, 0),
    # 0xFC = 'AKQJ109'
    (6, 0.4, 6, 6),
    # 0xFD = 'AKQJ1097'
    (7, 0, 7, 0),
    # 0xFE = 'AKQJ1098'
    (7, 0, 7, 6),
    # 0xFF = 'AKQJ10987'
    (8, 0, 8, 0)
)

NO_TRUMP_GAME = 4
MISERE_GAME = 5
LONG_TRUMP = 4
INITIAL_HAND_SIZE = 10
INITIAL_ORDERS_SIZE = 10
NUMBER_OF_HANDS = 3

PLAYER_NAMES = (u'Академик Новицкий', u'Академик Тархачёв')


class Suit(object):
    def __init__(self, code=0):
        self.code = code

    def __iadd__(self, rank_index):
        if 0 <= rank_index < SUIT_SIZE:
            if self.code & (1 << rank_index):
                raise IndexError
            self.code += (1 << rank_index)
            return self
        else:
            raise IndexError

    def __isub__(self, rank_index):
        if 0 <= rank_index < SUIT_SIZE:
            if not self.code & (1 << rank_index):
                raise IndexError
            self.code -= (1 << rank_index)
            return self
        else:
            raise IndexError

    def __str__(self):
        if not self.code:
            return '-'
        else:
            result = ''
            for i in xrange(SUIT_SIZE - 1, -1, -1):
                if self.code & (1 << i):
                    result += CARD_RANKS[i]
            return result

    def clear(self):
        self.code = 0

    def get_suit_cardinality(self):
        result = 0
        for i in xrange(SUIT_SIZE):
            if self.code & (1 << i):
                result += 1
        return result

    def get_indices(self):
        result = []
        for i in xrange(SUIT_SIZE):
            if self.code & (1 << i):
                result.append(i)
        return result

    def estimate_tricks(self, is_trump, trump_size=0, no_trump_game=False, misere=False):
        if misere:
            return TRICKS_ESTIMATIONS[self.code][MISERE]
        if no_trump_game:
            return TRICKS_ESTIMATIONS[self.code][NON_TRUMP_WITH_SHORT_TRUMP]
        if is_trump:
            return TRICKS_ESTIMATIONS[self.code][TRUMP]
        if trump_size < LONG_TRUMP:
            return TRICKS_ESTIMATIONS[self.code][NON_TRUMP_WITH_SHORT_TRUMP]
        else:
            return TRICKS_ESTIMATIONS[self.code][NON_TRUMP_WITH_LONG_TRUMP]


class Hand(object):
    def __init__(self, index_list=(), suits_code=0, horizontal=True):
        self.index_list = []
        self.horizontal = horizontal
        self.suits = [Suit(), Suit(), Suit(), Suit()]
        self.tricks_taken = 0
        if index_list:
            self.reset_by_list(index_list)
        elif suits_code:
            self.reset_by_code(suits_code)
        else:
            self.get_random_hand()

    def __iadd__(self, card_index):
        if 0 <= card_index < CARDS_SIZE:
            if card_index in self.index_list:
                raise IndexError
            else:
                self.index_list.append(card_index)
            self.suits[card_index / SUIT_SIZE] += (card_index % SUIT_SIZE)
            return self
        else:
            raise IndexError

    def __isub__(self, card_index):
        if 0 <= card_index < CARDS_SIZE:
            if card_index not in self.index_list:
                raise IndexError
            else:
                self.index_list.remove(card_index)
            self.suits[card_index / SUIT_SIZE] -= (card_index % SUIT_SIZE)
            return self
        else:
            raise IndexError

    def get_index_list_from_suit_codes(self):
        result = []
        for suit_index in xrange(SUITS_NUMBER):
            for rank_index in xrange(SUIT_SIZE):
                if self.suits[suit_index].code & (1 << rank_index):
                    result.append(SUIT_SIZE * suit_index + rank_index)
        return result

    def reset_by_list(self, index_list):
        self.tricks_taken = 0
        self.index_list = list(index_list)
        for index in index_list:
            if 0 <= index < CARDS_SIZE:
                self.suits[index / SUIT_SIZE] += (index % SUIT_SIZE)
            else:
                raise IndexError

    def reset_by_code(self, suits_code):
        for i in xrange(SUITS_NUMBER):
            self.suits[SUITS_NUMBER - i - 1].code = suits_code % 0x100
            suits_code >>= SUIT_SIZE
        self.index_list = self.get_index_list_from_suit_codes()

    def get_random_hand(self, size=INITIAL_HAND_SIZE):
        self.clear()
        card_indices = sample(xrange(CARDS_SIZE), size)
        for i in card_indices:
            self.__iadd__(i)

    def __str__(self):
        result = []
        for i in xrange(SUITS_NUMBER):
            result.append(SUIT_SYMBOLS[i] + ' ' + str(self.suits[i]))
        delimiter = u' ' if self.horizontal else u'\n'
        return delimiter.join(result).encode('utf-8')

    def clear(self):
        for suit in self.suits:
            suit.clear()

    def get_estimations(self):
        cardinalities = [suit.get_suit_cardinality() for suit in self.suits]
        result = []
        for trump in xrange(MISERE_GAME + 1):
            tricks_estimation = 0
            for suit in xrange(SUITS_NUMBER):
                est = self.suits[suit].\
                    estimate_tricks(trump == suit,
                                    cardinalities[trump] if trump < NO_TRUMP_GAME else 0,
                                    trump >= NO_TRUMP_GAME,
                                    trump == MISERE_GAME)
                tricks_estimation += est
            result.append(tricks_estimation)
        return result

    def print_tricks_estimations(self):
        estimations = self.get_estimations()
        for i in xrange(len(estimations)):
            print SUIT_SYMBOLS[i] + " : " + str(estimations[i]),

    def get_indices(self, suit):
        result = self.suits[suit].get_indices()
        for i in range(len(result)):
            result[i] += SUIT_SIZE * suit
        return result

    def get_available_indices(self, suit, trump):
        suit_indices = self.suits[suit].get_indices()
        if suit_indices:
            for i in xrange(len(suit_indices)):
                suit_indices[i] += (SUIT_SIZE * suit)
            return suit_indices
        if 0 <= trump < NO_TRUMP_GAME:
            trump_indices = self.suits[trump].get_indices()
            if trump_indices:
                for i in xrange(len(trump_indices)):
                    trump_indices[i] += (SUIT_SIZE * trump)
                return trump_indices
        return self.index_list


class Trick(object):
    def __init__(self, trump=-1):
        self.cards = []
        self.taker = -1
        self.first_hand = 0
        self.trump = trump

    def reset(self, first_hand):
        self.cards = []
        self.taker = -1
        self.first_hand = first_hand

    def less(self, first, second):
        if (first / SUIT_SIZE) == (second / SUIT_SIZE):
            return first < second
        else:
            return (second / SUIT_SIZE) == self.trump

    def __iadd__(self, card_index):
        self.cards.append(card_index)
        if len(self.cards) == 1:
            self.taker = 0
        else:
            if self.less(self.cards[self.taker], card_index):
                self.taker = len(self.cards) - 1
        return self

    def get_taker(self):
        return (self.first_hand + self.taker) % NUMBER_OF_HANDS


class Deal(object):
    def __init__(self, index_list, order=True):
        if len(index_list) != CARDS_SIZE:
            raise IndexError
        self.widow = index_list[3*INITIAL_HAND_SIZE:]
        self.order = order
        self.hands = [Hand(index_list[:INITIAL_HAND_SIZE], horizontal=False),
                      Hand(index_list[INITIAL_HAND_SIZE:2*INITIAL_HAND_SIZE]),
                      Hand(index_list[2*INITIAL_HAND_SIZE:3*INITIAL_HAND_SIZE])]
        self.trump = -1
        self.trick_order = -1
        self.hand_to_move = 0
        self.current_suit = -1
        self.trick = Trick()

    def reset(self, index_list):
        self.trump = -1
        self.trick_order = -1
        self.hand_to_move = 0
        self.current_suit = -1
        if len(index_list) != CARDS_SIZE:
            raise IndexError
        self.widow = index_list[3 * INITIAL_HAND_SIZE:]
        for i in range(NUMBER_OF_HANDS):
            self.hands[i].clear()
            self.hands[i].reset_by_list(index_list[INITIAL_HAND_SIZE*i:INITIAL_HAND_SIZE*(i+1)])

    def set_game_type(self, order, trump=-1):
        if order < MISERE_GAME:
            # ordered 6
            self.trump = order
            self.trick_order = 6
        elif order == MISERE_GAME:
            # ordered misere
            self.trick_order = 0
        else:
            # ordered 7, 8, 9, or 10
            self.trick_order = order + 1
            self.trump = trump
        self.trick.trump = self.trump

    def take_widow(self):
        self.hands[1] += self.widow[0]
        self.hands[1] += self.widow[1]
        self.widow = []

    def drop_cards(self, dropped_cards):
        self.hands[1] -= dropped_cards[0]
        self.hands[1] -= dropped_cards[1]
        self.widow = dropped_cards

    def get_available_moves(self, stage):
        if stage % NUMBER_OF_HANDS:
            return self.hands[self.hand_to_move % NUMBER_OF_HANDS].\
                get_available_indices(self.current_suit, self.trump)
        else:
            return self.hands[self.hand_to_move % NUMBER_OF_HANDS].index_list

    def make_move(self, stage, card_index):
        self.hands[self.hand_to_move] -= card_index
        if not (stage % NUMBER_OF_HANDS):
            self.trick.reset(self.hand_to_move)
            self.current_suit = card_index / SUIT_SIZE
        self.trick += card_index
        if stage % NUMBER_OF_HANDS != NUMBER_OF_HANDS - 1:
            self.hand_to_move = (self.hand_to_move + 1) % NUMBER_OF_HANDS
        else:
            self.hand_to_move = self.trick.get_taker()
            self.hands[self.hand_to_move].tricks_taken += 1

    def get_game_type(self):
        if self.trick_order == -1:
            return u'???'
        elif self.trick_order == 0:
            return u'МИЗЕР'
        else:
            result = str(self.trick_order)
            if self.trump != -1:
                result += SUIT_SYMBOLS[self.trump]
            return result

    def __str__(self):
        if self.order:
            result = 5 * ' ' + PLAYER_NAMES[1] + '\n' + 5 * ' ' +\
                     str(self.hands[1]).decode('utf-8') + '\n'
            west = str(self.hands[0]).decode('utf-8').split('\n')
            result += (west[0] + '\n')
            if self.widow:
                result += (west[1] + 20 * ' ' + SUIT_SYMBOLS[self.widow[0] / SUIT_SIZE] + ' ' +
                           CARD_RANKS[self.widow[0] % SUIT_SIZE] + 3 * ' ' + '\n')
                result += (west[2] + 20 * ' ' + SUIT_SYMBOLS[self.widow[1] / SUIT_SIZE] + ' ' +
                           CARD_RANKS[self.widow[1] % SUIT_SIZE] + 3 * ' ' + '\n')
            else:
                result += (west[1] + '\n')
                result += (west[2] + '\n')
            result += (west[3] + '\n')
            result += (5 * ' ' + PLAYER_NAMES[0] + '\n' + 5 * ' ' +
                       str(self.hands[2]).decode('utf-8'))
        else:
            result = 5 * ' ' + PLAYER_NAMES[1] + '\n' + 5 * ' ' +\
                     str(self.hands[2]).decode('utf-8') + '\n'
            east = str(self.hands[0]).decode('utf-8').split('\n')
            result += (25 * ' ' + east[0] + '\n')
            if self.widow:
                result += (2 * ' ' + SUIT_SYMBOLS[self.widow[0] / SUIT_SIZE] + ' '
                           + CARD_RANKS[self.widow[0] % SUIT_SIZE] + 20 * ' ' + east[1] + '\n')
                result += (2 * ' ' + SUIT_SYMBOLS[self.widow[1] / SUIT_SIZE] + ' '
                           + CARD_RANKS[self.widow[1] % SUIT_SIZE] + 20 * ' ' + east[2] + '\n')
            else:
                result += (25 * ' ' + east[1] + '\n')
                result += (25 * ' ' + east[2] + '\n')
            result += (25 * ' ' + east[3] + '\n')
            result += (5 * ' ' + PLAYER_NAMES[0] + '\n' + 5 * ' ' +
                       str(self.hands[1]).decode('utf-8'))
        return result.encode('utf-8')

    def widow_str(self):
        return SUIT_SYMBOLS[self.widow[0] / SUIT_SIZE] + ' ' +\
               CARD_RANKS[self.widow[0] % SUIT_SIZE] + '\n' + \
               SUIT_SYMBOLS[self.widow[1] / SUIT_SIZE] + ' ' +\
               CARD_RANKS[self.widow[1] % SUIT_SIZE]


class Player(object):
    def __init__(self, smartness):
        self.smartness = smartness
        self.tricks_bound = INITIAL_TRICKS_BOUND
        self.misere_bound = INITIAL_MISERE_BOUND
        self.score = 0

    def reset(self):
        self.tricks_bound = INITIAL_TRICKS_BOUND
        self.misere_bound = INITIAL_MISERE_BOUND
        self.score = 0

    @staticmethod
    def get_available_suits_for_order(order_list):
        for i in xrange(MISERE_GAME + 1, INITIAL_ORDERS_SIZE):
            if i in order_list:
                return list(xrange(SUITS_NUMBER + 1))
        result = []
        for suit in xrange(SUITS_NUMBER + 1):
            if suit in order_list:
                result.append(suit)
        return result

    def make_order(self, trick_estimations, order_list):
        if not (self.smartness & 1) or (len(order_list) == 1):
            return choice(order_list)
        available_suits = self.get_available_suits_for_order(order_list)
        misere_estimation = trick_estimations[-1]
        max_trick_estimation = 0
        for suit in xrange(len(trick_estimations)):
            if suit in available_suits and max_trick_estimation < trick_estimations[suit]:
                max_trick_estimation = trick_estimations[suit]

        available_big_games = []
        available_suits_for_six = []
        for order in order_list:
            if order > MISERE_GAME:
                available_big_games.append(order)
            if order < MISERE_GAME:
                available_suits_for_six.append(order)

        if MISERE_GAME in order_list and ((misere_estimation < self.misere_bound and
                                          max_trick_estimation < self.tricks_bound) or
                                          max_trick_estimation + 2 < self.tricks_bound):
            self.tricks_bound -= TRICKS_STEP
            return MISERE_GAME
        elif not available_suits_for_six or \
                (max_trick_estimation > self.tricks_bound and available_big_games):
            if len(available_big_games) == 1:
                order = available_big_games[0]
            elif len(available_big_games) == 2:
                if max_trick_estimation > self.tricks_bound + 3:
                    order = available_big_games[1]
                else:
                    order = available_big_games[0]
            elif len(available_big_games) == 3:
                if max_trick_estimation > self.tricks_bound + 4:
                    order = available_big_games[2]
                elif max_trick_estimation > self.tricks_bound + 2:
                    order = available_big_games[1]
                else:
                    order = available_big_games[0]
            else:
                if max_trick_estimation > self.tricks_bound + 3:
                    order = available_big_games[3]
                elif max_trick_estimation > self.tricks_bound + 2:
                    order = available_big_games[2]
                elif max_trick_estimation > self.tricks_bound + 1:
                    order = available_big_games[1]
                else:
                    order = available_big_games[0]
            if order not in order_list:
                raise ValueError
            self.tricks_bound += (order - 6) * TRICKS_STEP
            self.misere_bound += MISERE_STEP
            return order
        else:
            self.tricks_bound -= TRICKS_STEP
            self.misere_bound += MISERE_STEP
            if len(available_suits_for_six) == 1:
                return available_suits_for_six[0]
            max_trick_estimation = 0
            best_suit = choice(available_suits_for_six)
            for suit in available_suits_for_six:
                if trick_estimations[suit] > max_trick_estimation:
                    max_trick_estimation = trick_estimations[suit]
                    best_suit = suit
            if best_suit not in order_list:
                print order_list
                print available_suits_for_six
                print available_big_games
                raise ValueError
            return best_suit

    def choose_trump(self, hand):
        if not (self.smartness & 1):
            return randint(0, SUITS_NUMBER)
        cardinalities = [suit.get_suit_cardinality() for suit in hand.suits]
        max_suits = np.argwhere(cardinalities == np.max(cardinalities)).ravel().tolist()
        if len(max_suits) == 1:
            return max_suits[0]
        elif len(max_suits) == 2:
            if hand.suits[max_suits[0]].code < hand.suits[max_suits[1]].code:
                return max_suits[0]
            else:
                return max_suits[1]
        else:
            codes = [(i, hand.suits[i].code) for i in max_suits]
            codes = sorted(codes, key=lambda x: (x[1], x[0]))
            return codes[1][0]

    def drop_cards(self, hand, trump, mizere=False):
        if not self.smartness & 2:
            return sample(hand.index_list, 2)

        trick_estimation = 100 if mizere else 0
        list_index = copy(hand.index_list)
        for i in hand.index_list:
            if i / SUIT_SIZE == trump:
                list_index.remove(i)
        best_drop = sample(list_index, 2)
        for i in xrange(len(list_index)):
            if list_index[i] / SUIT_SIZE == trump:
                continue
            hand -= list_index[i]
            for j in xrange(i + 1, len(list_index)):
                if list_index[j] / SUIT_SIZE == trump:
                    continue
                hand -= list_index[j]
                trick_estimations_after_drop = hand.get_estimations()
                if mizere:
                    if trick_estimations_after_drop[MISERE_GAME] < trick_estimation:
                        trick_estimation = trick_estimations_after_drop[MISERE_GAME]
                        best_drop = [list_index[i], list_index[j]]
                else:
                    if trick_estimations_after_drop[trump] > trick_estimation:
                        trick_estimation = trick_estimations_after_drop[trump]
                        best_drop = [list_index[i], list_index[j]]
                hand += list_index[j]
            hand += list_index[i]
        return best_drop

    def make_move(self, index_list):
        return choice(index_list)


class Game(object):
    def __init__(self, first_player, second_player, order=False, verbose=0):
        self.players = [Player(first_player), Player(second_player)]
        self.order = order
        self.verbose = verbose
        self.deal = Deal(range(CARDS_SIZE), self.order)
        self.order_lists = [list(range(INITIAL_ORDERS_SIZE)), list(range(INITIAL_ORDERS_SIZE))]

    @staticmethod
    def get_score_of_ordering_player(game_type, tricks):
        if game_type < 5:
            if tricks < 6:
                return -48 + 32 * (tricks - 5)
            else:
                return 4 + 14 * (tricks - 6)
        elif game_type == 5:
            if tricks:
                return -100 * tricks
            else:
                return 100
        elif game_type == 6:
            if tricks < 7:
                return -120 + 64 * (tricks - 6)
            elif tricks == 7:
                return 16
            elif tricks == 8:
                return 24
            elif tricks == 9:
                return 52
            elif tricks == 10:
                return 80
            else:
                raise ValueError
        elif game_type == 7:
            if tricks < 8:
                return -120 + 96 * (tricks - 7)
            elif tricks == 8:
                return 36
            elif tricks == 9:
                return 48
            elif tricks == 10:
                return 90
            else:
                raise ValueError
        elif game_type == 8:
            if tricks < 9:
                return -144 + 128 * (tricks - 8)
            elif tricks == 9:
                return 64
            elif tricks == 10:
                return 120
            else:
                raise ValueError
        elif game_type == 9:
            if tricks < 10:
                return -160 * (10 - tricks)
            elif tricks == 10:
                return 150
            else:
                raise ValueError
        else:
            raise ValueError

    def make_random_deal(self):
        random_perm = range(CARDS_SIZE)
        shuffle(random_perm)
        self.deal.reset(random_perm)

    def get_order_list(self):
        result = []
        for order in self.order_lists[self.order]:
            if order < MISERE_GAME:
                result.append(str(6) + SUIT_SYMBOLS[order])
            elif order == MISERE_GAME:
                result.append(SUIT_SYMBOLS[order])
            else:
                result.append(str(order + 1))
        return ' '.join(result)

    def play_deal(self):
        self.make_random_deal()
        if self.verbose:
            print self.deal

        if self.verbose:
            print u"Доступные заказы игрока {}:".format(PLAYER_NAMES[self.order])
            print self.get_order_list()

        order_index = self.players[self.order].make_order(self.deal.hands[1].get_estimations(),
                                                          self.order_lists[self.order])
        self.order_lists[self.order].remove(order_index)
        self.deal.set_game_type(order_index)
        self.deal.take_widow()
        if self.verbose:
            print PLAYER_NAMES[int(self.order)] + u' заказал ' +\
                  self.deal.get_game_type() + u' и взял прикуп'
        if self.verbose == 2:
            print self.deal

        trump_index = order_index
        if order_index > MISERE_GAME:
            trump_index = self.players[self.order].choose_trump(self.deal.hands[1])
        self.deal.set_game_type(order_index, trump_index)
        if self.verbose:
            print PLAYER_NAMES[int(self.order)] + u' играет ' + self.deal.get_game_type()

        dropped_cards = self.players[self.order].drop_cards(self.deal.hands[1],
                                                            trump_index,
                                                            order_index == MISERE_GAME)
        self.deal.drop_cards(dropped_cards)
        if self.verbose:
            print PLAYER_NAMES[int(self.order)] + u' сделал снос '
        if self.verbose:
            print self.deal

        for stage in range(NUMBER_OF_HANDS * INITIAL_HAND_SIZE):
            available_moves = self.deal.get_available_moves(stage)
            player_index = (self.deal.hand_to_move + self.order) % 2
            move = self.players[player_index].make_move(available_moves)
            if self.verbose == 2:
                print PLAYER_NAMES[1 - player_index] + u' сделал ход ' +\
                     SUIT_SYMBOLS[move / SUIT_SIZE] + ' ' + CARD_RANKS[move % SUIT_SIZE]
            self.deal.make_move(stage, move)
            if self.verbose == 2 and stage % NUMBER_OF_HANDS == 2:
                print u'Взятку взяла рука {}'.format(self.deal.hand_to_move + 1)
            if (self.verbose == 2 and stage % NUMBER_OF_HANDS == 2) or\
                    (self.verbose and stage == NUMBER_OF_HANDS*INITIAL_HAND_SIZE - 1):
                print u"Текущее распредление взяток:" \
                      u" {}:{}:{}".format(self.deal.hands[0].tricks_taken,
                                          self.deal.hands[1].tricks_taken,
                                          self.deal.hands[2].tricks_taken)
            if self.verbose == 2 and stage % NUMBER_OF_HANDS == 2 \
                    and stage < NUMBER_OF_HANDS*INITIAL_HAND_SIZE - 1:
                print self.deal

        tricks_of_ordering_player = self.deal.hands[1].tricks_taken
        score_of_ordering_player = self.get_score_of_ordering_player(order_index,
                                                                     tricks_of_ordering_player)
        current_score = -score_of_ordering_player if self.order else score_of_ordering_player
        if self.verbose:
            print u"Счёт текущей сдачи: {}:{}".format(current_score, -current_score)
        self.players[self.order].score += score_of_ordering_player
        self.players[not self.order].score -= score_of_ordering_player
        if self.verbose:
            print u"Общий счёт в партии: {}: {}, {}: {}".\
                format(PLAYER_NAMES[0], self.players[0].score,
                       PLAYER_NAMES[1], self.players[1].score)

    def play_game(self):
        for i in xrange(2 * INITIAL_ORDERS_SIZE):
            self.play_deal()
            self.order = not self.order
            self.deal.order = self.order
        if self.verbose:
            print u"Счёт партии: {}: {}, {}: {}".\
                format(PLAYER_NAMES[0], self.players[0].score,
                       PLAYER_NAMES[1], self.players[1].score)
        return self.players[0].score

    def __str__(self):
        return str(self.deal)


def test(attempts=100, first_player=0, second_player=0, verbose=0):
    array = []
    for i in xrange(attempts):
        game = Game(first_player, second_player, verbose=verbose)
        array.append(game.play_game())

    scores = np.array(array)
    print u"Средний выигрыш игрока {}: {}".format(PLAYER_NAMES[0], np.mean(array))
    print u"Соотношение побед и поражений игрока {}: {}:{}".\
        format(PLAYER_NAMES[0], (scores > 0).sum(), (scores < 0).sum())


def execute_command(args):
    if args.command == 'hand':
        code = int(args.code, 16)
        if code:
            hand = Hand(suits_code=code)
        else:
            hand = Hand()
        print u"Что нам раздали: ", hand
        print u"Ожидание числа взяток на различных играх:"
        hand.print_tricks_estimations()
    elif args.command == 'misere':
        print u"В поисках мизера..."
        counter = 0
        for i in xrange(args.attempts):
            hand = Hand()
            trick_estimations = hand.get_estimations()
            if trick_estimations[-1] < args.bound:
                counter += 1
                print hand
        print u"Найдено {} околомизерных сдач в {} случайных сдачах," \
              u" вероятность мизера {}".\
            format(counter, args.attempts, float(counter) / float(args.attempts))
    elif args.command == 'suit':
        print u"В поисках игры " + SUIT_SYMBOLS[SUIT_NAMES[args.trump]] + ' ...'
        counter = 0
        for i in xrange(args.attempts):
            hand = Hand()
            trick_estimations = hand.get_estimations()
            if trick_estimations[SUIT_NAMES[args.trump]] > args.bound:
                counter += 1
                print hand
        print u"Найдено {} пригодных игр в {} случайных сдачах," \
              u" вероятность игры {}".\
            format(counter, args.attempts, float(counter) / float(args.attempts))
    elif args.command == 'deal':
        game = Game(args.first, args.second, bool(args.order % 2), args.verbose)
        game.play_deal()
    elif args.command == 'game':
        for i in xrange(args.number):
            game = Game(args.first, args.second, bool(args.order % 2), args.verbose)
            game.play_game()


def main():
    parser = argparse.ArgumentParser(description='Arguments for preferance.')

    subparsers = parser.add_subparsers(dest='command')

    parser_hand = subparsers.add_parser('hand', help='Generate hand by code or random.')
    parser_hand.add_argument('-c', '--code', type=str,
                             default='0', help='Hex code, default - random hand.')

    parser_misere = subparsers.add_parser('misere',
                                          help='Searching random hands for playing misere.')
    parser_misere.add_argument('-a', '--attempts', type=int, default=100, help='Attempts size.')
    parser_misere.add_argument('-b', '--bound', type=float, default=1.0,
                               help='Bound for misere coefficient.')

    parser_suit = subparsers.add_parser('suit',
                                        help='Searching random hands for playing certain suit.')
    parser_suit.add_argument('trump', type=str, help='Trump suit.')
    parser_suit.add_argument('-a', '--attempts', type=int, default=100, help='Attempts size.')
    parser_suit.add_argument('-b', '--bound', type=float, default=5.0,
                             help='Bound for playing game.')

    parser_deal = subparsers.add_parser('deal', help='Play random deal.')
    parser_deal.add_argument('-o', '--order', type=int, default=0, help='Player who makes order.')
    parser_deal.add_argument('-f', '--first', type=int, default=0, help='First player smartness.')
    parser_deal.add_argument('-s', '--second', type=int, default=0,
                             help='Second player smartness.')
    parser_deal.add_argument('-v', '--verbose', type=int, default=0, help='Verbose level.')

    parser_game = subparsers.add_parser('game', help='Play bandit20.')
    parser_game.add_argument('-o', '--order', type=int, default=0, help='Player who orders first.')
    parser_game.add_argument('-f', '--first', type=int, default=0, help='First player smartness.')
    parser_game.add_argument('-s', '--second', type=int, default=0,
                             help='Second player smartness.')
    parser_game.add_argument('-v', '--verbose', type=int, default=0, help='Verbose level.')
    parser_game.add_argument('-n', '--number', type=int, default=1, help='Number of games.')

    args = parser.parse_args()
    execute_command(args)


if __name__ == '__main__':
    main()
