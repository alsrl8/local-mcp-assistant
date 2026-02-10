# mcp-assistant

Go로 작성된 경량 MCP(Model Context Protocol) 서버입니다. WSL 환경에서 동작하며, Claude Desktop 등 MCP 클라이언트에 로컬 파일시스템 접근 도구를 제공합니다. Windows 경로를 자동으로 WSL 경로로 변환하여 별도 설정 없이 파일 탐색과 검색이 가능합니다.

## 아키텍처

![architecture](/docs/architecture_v1_0_0.png)

## 제공 도구

| 도구 | 설명 |
|------|------|
| `list_dir` | 디렉토리 내 파일/폴더 목록 조회 (이름, 크기, 타입) |
| `read_file` | 파일 내용 읽기 (512KB 초과 시 자동 truncate) |
| `grep` | 정규식 기반 파일/디렉토리 내 텍스트 검색 (context line 지원) |

모든 도구는 Windows 경로(`C:\Users\...`)를 입력받아 WSL 경로(`/mnt/c/Users/...`)로 자동 변환합니다.

## 기술 스택

- **Go 1.25** — [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- **Stdio Transport** — MCP 통신
- **Viper** — 설정 관리
- **WSL** — 런타임 환경

## 시작하기

### 사전 요구사항

- WSL (Ubuntu 권장)
- Go 1.25+

### 빌드 및 설치

```bash
# 빌드
make build

# ~/.local/bin에 설치
make install

# 빌드 + 설치
make deploy
```

### Claude Desktop 설정

`claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "assistant": {
      "command": "wsl.exe",
      "args": ["/home/<username>/.local/bin/mcp-assistant"],
      "env": {
        "MCP_ASSISTANT_CONFIG": "/path/to/config.yaml"
      }
    }
  }
}
```

## 프로젝트 구조

```
├── main.go          # 진입점 — 서버 초기화 및 stdio 트랜스포트 실행
├── tool/            # 도구 구현 (도구당 파일 1개)
│   └── init.go      # 도구 등록
├── schema/          # 각 도구의 입출력 타입 정의
├── config/          # Viper 기반 설정 (DB 설정 등 확장용)
├── utils/           # 경로 변환 (Windows → WSL)
├── example/         # MCP 도구 구현 예제 (SayHi)
└── docs/            # 설정 가이드 및 사용 예시
```

## 라이선스

MIT
