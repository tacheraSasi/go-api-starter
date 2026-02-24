import { z } from 'zod';

const isEmailValid = (email: string) => z.string().email().safeParse(email).success;

const isPositive = (num: number) => z.number().positive().safeParse(num).success;

const isNegative = (num: number) => z.number().negative().safeParse(num).success;

export { isEmailValid, isPositive, isNegative };
