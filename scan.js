#!/usr/bin/env node

const args = process.argv.slice(2);
const { mainModule } = require('process');
const util = require('util');
const exec = util.promisify(require('child_process').exec);

// \x1b[0m：重置所有属性，恢复终端默认颜色。
//   \x1b[30m - \x1b[37m：设置前景色，范围是 30 到 37，分别对应黑、红、绿、黄、蓝、紫、青、白。
//   \x1b[40m - \x1b[47m：设置背景色，范围是 40 到 47，分别对应黑、红、绿、黄、蓝、紫、青、白。
//   \x1b[1m：设置加粗。
//   \x1b[4m：设置下划线。
//   \x1b[7m：设置反显。

function yellow(str) {
  return `\x1b[33m${str}\x1b[0m`;
}
function red(str) {
  return `\x1b[31m${str}\x1b[0m`;
}

async function main() {

  const coderIp = '172.16.158.240'
  const ports = args.length ? args : (await scanPorts()).filter(port => port >= 3000 && port < 10000);

  const {userWorksapceName, user, workspace} = getWorkspaceName();

  function nip(port, ip) {
    return `https://${port}-${userWorksapceName}-${toHex(ip)}.nip.io`;
  }
  
  function gltnuxt(port) {
    return `https://${port}-${userWorksapceName}.gltnuxt.top`;
  }

  function chanjet(port) {
    return `https://code.chanjet.com.cn/proxy/${user}/${workspace}/${port}`;
  }

  console.log(`
  port 3000 6000   返回指定端口的可访问域名
  port             自动扫描3000到10000可用端口，并返回可访问域名

  ${yellow('警告: 如提示证书不安全,请全局安装@chanjet/cjet-proxy, cjet代理程序启动时会自动安装证书')}
`)

  ports.forEach((port) => {
    console.log(`
  ${port}  ->  ${red(gltnuxt(port))}`)
  })

  console.log('')
}

main();



// env

function getWorkspaceName() {
  const { FRP_NAME, GIT_AUTHOR_NAME: user, HOSTNAME: workspace } = process.env;
  const userWorksapceName = FRP_NAME || GIT_AUTHOR_NAME + '-' + HOSTNAME;
  return {userWorksapceName, user, workspace }
}

// 10进制 to hex
function toHex(str) {
  return str.split('.').map((item) => {
    // console.log(('0' + parseInt(item, 10).toString(16)).slice(-2))
    return ('0' + parseInt(item, 10).toString(16)).slice(-2);
  }).join('');
}


async function scanPorts() {

  const { stdout } = await exec('lsof -i -P -n | grep LISTEN');

  const lines = stdout.trim().split('\n');
  const ports = lines.map(line => {
    const parts = line.trim().split(/\s+/);
    const address = parts[8];
    const port = address.split(':')[1];
    return parseInt(port);
  });

  return ports
}