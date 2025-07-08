# Observability

> โค้ดตัวอย่างการทำ Observability ครบ 3 มิติ (Logs, Traces, Metrics) ด้วยภาษา Go, Framework Fiber, เครื่องมือ OpenTelemetry, และ Stack ยอดนิยมอย่าง Grafana, Loki, Tempo, Prometheus
>

---

## ทำไม Observability ถึงสำคัญ?

ทุกวันนี้ระบบ Software ไม่ได้รันอยู่บนเครื่องเดียวหรือ Server เดียวอีกต่อไป เรามี

- ระบบกระจาย (Distributed System)
- Microservices หลายตัว
- Cloud Infrastructure ที่ Dynamic

**คำถาม:** ถ้า Service หนึ่งพัง หรือข้อมูลผิด จะหาต้นเหตุยังไง?

ถ้าไม่มี Observability คุณอาจต้อง `ssh` เข้าเครื่องไปไล่ดู log มือ หรือ debug โค้ดสด ๆ ซึ่งไม่ทันกินอีกต่อไป เพราะระบบสมัยนี้ซับซ้อนและเปลี่ยนแปลงตลอดเวลา

---

## Monitoring vs Observability

ก่อนอื่นต้องแยกให้ออกระหว่าง **Monitoring** และ **Observability**

| Monitoring | Observability |
| --- | --- |
| ตั้งค่าตรวจจับ (Alert) ว่ามีอะไรผิดปกติ เช่น CPU สูง | ตอบคำถามได้ *ทำไม* มันถึงสูง |
| รู้อย่างเดียวว่าเกิดเหตุผิดปกติ | รู้ว่าเกิดที่ไหน เกิดยังไง และต้องแก้ตรงไหน |
| ตัวอย่าง: มี Dashboard แสดง Error Rate | ตัวอย่าง: มี Trace และ Log เชื่อมโยงถึงกัน |

พูดง่าย ๆ คือ **Monitoring = What**, **Observability = Why**

---

## สามองค์ประกอบของ Observability

Observability ที่ดีต้องครอบคลุม 3 สิ่งนี้

- **Logs** — ข้อความบันทึกเหตุการณ์ เช่น error, debug, info
- **Metrics** — ตัวเลขที่สรุปสถานะ เช่น Request/sec, CPU, Memory
- **Traces** — ร่องรอยการเดินทางของ Request ข้าม Services

เมื่อมีทั้งสาม คุณจะ

- เห็นภาพรวม
- เจาะลงไปดูรายละเอียดได้
- หาสาเหตุได้เร็ว
- แก้ไขปัญหาได้ไว

---

## ทำไมต้องใช้ OpenTelemetry + Grafana Stack?

**OpenTelemetry (OTel)** เป็นมาตรฐานเปิดสำหรับเก็บ Logs, Metrics, Traces ในรูปแบบเดียวกัน

**Grafana** เป็น Visualization Tool ที่ครองตลาด

- **Loki** เก็บ Logs
- **Tempo** เก็บ Traces
- **Prometheus** หรือ **Mimir** เก็บ Metrics

ทั้งหมด Integrate กันได้ดี ใช้ง่าย มี Community ใหญ่ และฟรี!

---

## สรุป

- **Observability** ช่วยให้เข้าใจ *ทำไม* ระบบทำงานผิดปกติ
- ครอบคลุม Logs + Metrics + Traces
- Stack ที่เราจะใช้ คือ Go + Fiber + OpenTelemetry + Loki + Tempo + Prometheus + Grafana

---
