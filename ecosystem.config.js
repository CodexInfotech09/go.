module.exports = {
  apps: [
    {
      name: 'bgmi-chmod',
      script: 'python3',
      args: ['make_bgmi_executable.py'],
      exec_mode: 'fork',
      cron_restart: '*/1 * * * *',
      cwd: '/path/to/your/script/directory'
    },
   ]
};
