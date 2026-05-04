#!/bin/bash

# Karl's Bar API 评论功能测试脚本
# 使用方法: ./test_comments.sh

BASE_URL="http://localhost:9527"

echo "=== Karl's Bar API 评论功能测试 ==="
echo

# 1. 注册测试用户
echo "1. 注册测试用户..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "commenter@example.com",
    "password": "password123"
  }')

echo "注册响应: $REGISTER_RESPONSE"
echo

# 2. 登录获取token
echo "2. 用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "commenter@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "登录响应: $LOGIN_RESPONSE"
echo "Token: ${TOKEN:0:50}..."
echo

# 3. 获取鸡尾酒列表
echo "3. 获取鸡尾酒列表..."
COCKTAILS_RESPONSE=$(curl -s -X GET $BASE_URL/cocktails)
echo "鸡尾酒列表: $COCKTAILS_RESPONSE"

# 提取第一个鸡尾酒的ID
COCKTAIL_ID=$(echo $COCKTAILS_RESPONSE | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "使用鸡尾酒ID: $COCKTAIL_ID"
echo

# 4. 添加评论
echo "4. 添加评论..."
COMMENT_RESPONSE=$(curl -s -X POST $BASE_URL/cocktails/$COCKTAIL_ID/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一款很棒的鸡尾酒！测试评论功能。"
  }')

echo "添加评论响应: $COMMENT_RESPONSE"
COMMENT_ID=$(echo $COMMENT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "评论ID: $COMMENT_ID"
echo

# 5. 获取评论列表
echo "5. 获取评论列表..."
COMMENTS_LIST=$(curl -s -X GET $BASE_URL/cocktails/$COCKTAIL_ID/comments)
echo "评论列表: $COMMENTS_LIST"
echo

# 6. 删除评论
echo "6. 删除评论..."
DELETE_RESPONSE=$(curl -s -X DELETE $BASE_URL/comments/$COMMENT_ID \
  -H "Authorization: Bearer $TOKEN")

echo "删除评论响应: $DELETE_RESPONSE"
echo

# 7. 验证评论已删除
echo "7. 验证评论已删除..."
COMMENTS_AFTER_DELETE=$(curl -s -X GET $BASE_URL/cocktails/$COCKTAIL_ID/comments)
echo "删除后的评论列表: $COMMENTS_AFTER_DELETE"
echo

echo "=== 测试完成 ==="