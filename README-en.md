# GAPChain MVP — Open-Source Agri-Blockchain

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Fabric](https://img.shields.io/badge/Hyperledger%20Fabric-3.1.1-2C3E50.svg)](https://hyperledger.org)
[![Go Base](https://img.shields.io/badge/Golang-1.21-00ADD8.svg)](https://golang.org/)
[![Vue 3](https://img.shields.io/badge/Vue.js-3.x-4FC08D.svg)](https://vuejs.org/)

Welcome to **GAPChain**, an open-source agricultural traceability ecosystem. This projected distributed ledger is not crafted for generic marketing pitches—it's a technical foundation built to mathematically enforce transparency and protect the organic farming supply chain in Southeast Asia. 

**We actively need your contribution.** The GAPChain MVP is smoothly running locally orchestrated over a 4-Org architecture, 2 Secure Channels, and 3 Chaincodes. We are looking for brave community developers (Blockchain Engineers, Backend, DevOps, Flutter Devs) who share the vision to help scale this architecture up to a massive Production system!

---

## 1. Why Should You Contribute?

This is not another boilerplate CRUD web-app. By joining our maintainer ranks, you will solve rigorous enterprise-wide distributed ledger challenges:

1. **Precision Yield Tokenization:** Advanced chaincode dynamics ensure exactly 1 metric ton of recorded harvest translates to precisely 1 ton of trackable yield. Mathematically denying illicit QR copy printing.
2. **Immutability Audit Trails:** Exploiting the `GetHistoryForKey` CouchDB primitive returning the entire life cycles. Admin override of historical timestamps is cryptographically blocked.
3. **Event-driven Geo-Clone Shields:** We log metadata scans to alert end-users real-time whenever physical label duplicates are scanned across contradictory geo-coordinates simultaneously.
4. **On-Chain Label Burning:** Consumer confirmation "burns" the package's state, preventing the physical outer box from being repackaged by counterfeiters down the line.
5. **MSP Identity Policies:** Hardened API layers utilizing robust `contractapi` signatures ensuring solely specific governmental peers (BVTV) invoke quality certificates. 

> 🛠️ *"Every single Pull Request—be it a three-line CI/CD configuration tweak or a heavily optimized Flutter module—will be personally reviewed, merged, and fully credited in our documentation."*

---

## 2. Technical Roadmap & Help Wanted

- [x] **Phase 1: Local MVP (✅ Complete)** 
  - Stood up the 4-Organization local network communicating via modern gRPC Protocol `fabric-gateway-go` v1.x.
  - Operationalized three Core Go Smart Contracts: `lohang_cc` (Batches), `nhatky_cc` (Journals), `giao_dich_cc` (Trades).
  - Live Frontend Consumer web scanning portal using Vue 3.
- [ ] **Phase 2: Cloud Deployment (🔥 Next Priorities)**
  - Implement **CI/CD Pipelines** triggering automations via GitHub Actions.
  - Setup bullet-proof **Docker Compose Production** topologies tailored for Cloud AWS or bare-metal VPS hosting solutions.
  - Embed health telemetries leveraging **Prometheus + Grafana Dashboards**.
  - Execute strict test coverage paths aiming for `>70%` integration checkpoints (Ginkgo).
- [ ] **Phase 3: Real-World Pilot Integrations**
  - Consolidate **Offline-first Local Sync protocols** within our Flutter App logic—supporting rural farmers operating without 4G reception. 
  - Standardize logic hashing in-field camera outputs (SHA256) embedded with unchangeable GPS coords.
  - Attach a user-facing **Hyperledger Fabric Explorer**.

---

## 3. Architecture & Tech Stack

GAPChain utilizes a multi-layer distributed proxy architecture powered by Hyperledger Fabric to shield end-users.

### 3.1 Core Stack Elements
- **Blockchain Core**: Hyperledger Fabric v3.1.1.
- **Smart Contracts (Chaincode)**: Golang 1.19+ (leveraging `fabric-chaincode-go/v2`, `fabric-contract-api-go/v2`).
- **State Database**: CouchDB (JSON formatted rich queries).
- **Backend API Modules**: Golang 1.21+ (Gin framework, DI using `uber-go/fx`).
- **Mobile Ecosystem**: Flutter 3.x, Riverpod, sqflite.
- **Consumer Web App**: Vue 3, Vite.

### 3.2 Network Topology
Built over 4 isolated peer Organizations:
1. **PlatformOrg** (`platform.gapchain.vn`): Overarching admin logic + emergency suspensions.
2. **HTXNongSanOrg** (`htxnongsan.gapchain.vn`): Cooperative farmers managing crops.
3. **ChiCucBVTVOrg** (`chicucbvtv.gapchain.vn`): Official verification agency nodes.
4. **NPPXanhOrg** (`nppxanh.gapchain.vn`): Supermarket retail procurement chains.

Overhead is optimized across 2 channels: `nhatky-htx-channel` (Ops/Logs) & `giaodich-channel` (Sales/Trades).

---

## 4. Development Guides (Backend, Golang)

The Golang Backend enforces security and maps Web JWTs -> Fabric MSP logic securely. 

```text
gapchain/backend/
├── cmd/server/             # Primary runner loading uber-go/fx dependencies
├── internal/
│   ├── infrastructure/     # Drivers, Fabric Gateway interfaces
│   ├── repository/fabric/  # EvaluateTransaction (Read) & SubmitTransaction (Execute/Commit)
│   ├── middleware/         # Security mappings
│   └── model/              # HTTP Data Transfer Objects
```

---

## 5. Setup & Local Demo Bootstrapping

Requires Docker Engine + Docker Compose installed.

```bash
cd gapchain/chain-setup
export PATH=$PATH:$(pwd)/../bin

# 1. Boot up Network (Will run 9 containers under standard 1 Orderer, 4 Peers layout)
./scripts/setup-network.sh

# 2. Package and Install Go Chaincodes 
./scripts/deploy-chaincode.sh

# 3. Spin up Backend Middlewares
cd ../backend
go run cmd/server/main.go
# Available locally at: localhost:8080

# 4. Demo the Vite Consumer Portal
cd ../frontend-web
npm install && npm run dev
```

> **Cleanup Tip:** Completely tear down networks using `./scripts/cleanup-network.sh --remove-images`.

---

## 6. Open Code Conventions

1. Chaincode JSON structs marshalled into CouchDB stick strictly to `snake_case`. Standardly accessible CC methods map to upper CamelCase variants (`TaoLotHang`, `GhiNhatKy`).
2. Determinism mandates avoiding native `time.Now()` operations inside Blockchain bounds. You **must** utilize the Fabric internal equivalent `ctx.GetStub().GetTxTimestamp()` guaranteeing matching consensus times securely.
3. Feature development follows `feature/your-contribution` tags. Commit messages preferably reflect semantics standard. 

🌟 _Thank you globally distributed community for supporting true open traceability ecosystems!_
