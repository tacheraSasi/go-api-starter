import { z } from 'zod';

// t: TFunction<'errors'[], undefined>

const loginSchemaValidation = () => {
  return z.object({
    email: z
      .string()
      .email({ message: 'Email is invalid' })
      .refine((val) => val.length > 0, {
        message: 'Email required',
      }),
    password: z.string().refine((val) => val.length > 0, {
      message: 'Password Required',
    }),
  });
};

export { loginSchemaValidation };
