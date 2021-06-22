# 第一步：获取运行环境
FROM ubuntu:21.04   


# 第二步：拷贝可执行文件
COPY . /home/
RUN chmod 777 /home/main

# 第三步: 运行文件
CMD /home/env/main



