import pinoLogger from 'pino';

const logger = pinoLogger({
  transport: {
    options: {
      colorize: true,
    },
    target: 'pino-pretty',
  },
});

export default logger;
