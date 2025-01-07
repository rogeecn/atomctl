{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "target": {
          "limit": 100,
          "matchAny": false,
          "tags": [],
          "type": "dashboard"
        },
        "type": "dashboard"
      }
    ]
  },
  "description": "otelcol metrics dashboard",
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": 3,
  "links": [],
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 8,
      "panels": [],
      "title": "Process",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Otel Collector Instance",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 4,
        "w": 3,
        "x": 0,
        "y": 1
      },
      "id": 6,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "center",
        "orientation": "auto",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "count(count(otelcol_process_cpu_seconds_total{service_instance_id=~\".*\"}) by (service_instance_id))",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Instance",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "otelcol_process_cpu_seconds",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "s"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 4,
        "w": 3,
        "x": 3,
        "y": 1
      },
      "id": 24,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "avg(rate(otelcol_process_cpu_seconds_total{}[$__rate_interval])*100) by (instance)",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Cpu",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Memory Rss\navg(otelcol_process_memory_rss{}) by (instance)",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          },
          "unit": "bytes"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 4,
        "w": 3,
        "x": 6,
        "y": 1
      },
      "id": 38,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "avg(otelcol_process_memory_rss{}) by (instance)",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Memory",
      "type": "gauge"
    },
    {
      "description": "",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "gridPos": {
        "h": 4,
        "w": 15,
        "x": 9,
        "y": 1
      },
      "id": 32,
      "options": {
        "code": {
          "language": "plaintext",
          "showLineNumbers": false,
          "showMiniMap": false
        },
        "content": "## Opentelemetry Collector Data Ingress/Egress\n\n`service_version:` ${service_version}\n\n`opentelemetry collector:` contrib\n\n",
        "mode": "markdown"
      },
      "pluginVersion": "11.3.0",
      "title": "",
      "type": "text"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 5
      },
      "id": 10,
      "panels": [],
      "title": "Trace Pipeline",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "(avg(sum by(job) (rate(otelcol_exporter_sent_spans{}[$__range]))) / avg(sum by(job) (rate(otelcol_receiver_accepted_spans{}[$__range])))) ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "min": 0,
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "light-blue",
                "value": null
              },
              {
                "color": "semi-dark-red",
                "value": 0
              },
              {
                "color": "super-light-orange",
                "value": 0.4
              },
              {
                "color": "dark-blue",
                "value": 0.9
              },
              {
                "color": "super-light-orange",
                "value": 1.2
              },
              {
                "color": "dark-red",
                "value": 2.1
              }
            ]
          },
          "unit": "percentunit"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 19,
        "w": 3,
        "x": 0,
        "y": 6
      },
      "id": 55,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "vertical",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": false,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "avg(sum by(job) (rate(otelcol_exporter_sent_spans_total{}[$__range])))",
          "format": "time_series",
          "hide": true,
          "instant": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "export"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "avg(sum by(job) (rate(otelcol_receiver_accepted_spans_total{}[$__range])))",
          "format": "time_series",
          "hide": true,
          "legendFormat": "__auto",
          "range": true,
          "refId": "acc"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "(avg(sum by(job) (rate(otelcol_exporter_sent_spans_total{}[$__range]))) / avg(sum by(job) (rate(otelcol_receiver_accepted_spans_total{}[$__range])))) ",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Export Ratio",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "gridPos": {
        "h": 11,
        "w": 21,
        "x": 3,
        "y": 6
      },
      "id": 4,
      "options": {
        "edges": {},
        "nodes": {
          "mainStatUnit": "flops"
        }
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "label_join(label_join(\n(rate(otelcol_receiver_accepted_spans_total{}[$__interval]))\n, \"id\", \"\", \"transport\", \"receiver\")\n, \"title\", \"\", \"transport\", \"receiver\")\n\nor\n\nlabel_replace(label_replace(\nsum by(service_name) (rate(otelcol_receiver_accepted_spans_total{}[$__interval]))\n, \"id\", \"processor\", \"dummynode\", \"\")\n, \"title\", \"processor\", \"dummynode\", \"\")\n\nor\nlabel_replace(label_replace(\n(rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n, \"id\", \"$0\", \"processor\", \".*\")\n, \"title\", \"$0\", \"processor\", \".*\")\n\nor\nlabel_replace(label_replace(\nsum by(exporter) (rate(otelcol_exporter_sent_spans_total{}[$__interval]))\n, \"id\", \"exporter\", \"dummynode\", \"\")\n, \"title\", \"exporter\", \"dummynode\", \"\")\n        \nor\nlabel_replace(label_replace(\nsum by(exporter) (rate(otelcol_exporter_sent_spans_total{}[$__interval]))\n, \"id\", \"$0\", \"exporter\", \".*\")\n, \"title\", \"$0\", \"exporter\", \".*\")",
          "format": "table",
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "nodes"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "label_join(\nlabel_replace(label_join(\n(rate(otelcol_receiver_accepted_spans_total{}[$__interval]))\n\n ,\"source\",\"\",\"transport\",\"receiver\")\n,\"target\",\"processor\",\"\",\"\")\n,\"id\",\"-\",\"source\",\"target\")\n\n  or\n\n  label_join(\nlabel_replace(label_replace(\n  (rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n ,\"source\",\"processor\",\"\",\"\")\n,\"target\",\"$0\",\"processor\",\".*\")\n,\"id\",\"-\",\"source\",\"target\")\n\nor\n  label_join(\nlabel_replace(label_replace(\n    (rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n ,\"source\",\"$0\",\"processor\",\".*\")\n,\"target\",\"exporter\",\"\",\"\")\n,\"id\",\"-\",\"source\",\"target\")\n\nor\n  label_join(\nlabel_replace(label_replace(\n   (rate(otelcol_exporter_sent_spans_total{}[$__interval]))\n ,\"source\",\"exporter\",\"\",\"\")\n,\"target\",\"$0\",\"exporter\",\".*\")\n,\"id\",\"-\",\"source\",\"target\")\n\n",
          "format": "table",
          "hide": false,
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "edges"
        }
      ],
      "title": "",
      "type": "nodeGraph"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Spans Accepted by Receiver and Transport",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "noValue": "no data",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 3,
        "y": 17
      },
      "id": 12,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_accepted_spans_total{}[$__rate_interval])) by (receiver,transport)",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Accepted",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Spans Accepted ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 3,
        "x": 8,
        "y": 17
      },
      "id": 13,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "horizontal",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_accepted_spans_total{}[$__rate_interval])) ",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Batch Processed",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "#EAB839",
                "value": 1
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 11,
        "y": 17
      },
      "id": 15,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "text": {},
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_processor_batch_batch_send_size_sum{}[$__rate_interval]))  by (processor)",
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Batch",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 3,
        "x": 16,
        "y": 17
      },
      "id": 14,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_sent_spans_total{}[$__interval])) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Sent by Exporter",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 19,
        "y": 17
      },
      "id": 30,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_sent_spans_total{}[$__rate_interval])) by (exporter)",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Sent",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "noValue": "no data",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 5,
        "x": 3,
        "y": 22
      },
      "id": 17,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_refused_spans_total{}[$__rate_interval])) by (receiver,transport)",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Refused",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Spans Accepted ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 8,
        "y": 22
      },
      "id": 18,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_refused_spans_total{}[$__rate_interval])) ",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "otelcol_exporter_send_failed_spans",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 16,
        "y": 22
      },
      "id": 19,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_send_failed_spans_total{}[$__rate_interval])) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Sent by Exporter\notelcol_exporter_send_failed_spans",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 5,
        "x": 19,
        "y": 22
      },
      "id": 20,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_send_failed_spans_total{}[$__rate_interval])) by (exporter)",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Failed",
      "type": "stat"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 25
      },
      "id": 22,
      "panels": [],
      "title": "Metrics Pipeline",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "avg(sum by(job) (rate(otelcol_exporter_sent_metric_points{}[$__range]))) versus avg(sum by(job) (rate(otelcol_receiver_accepted_metric_points{}[$__range])))",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "min": 0,
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "light-blue",
                "value": null
              },
              {
                "color": "semi-dark-red",
                "value": 0
              },
              {
                "color": "super-light-orange",
                "value": 0.4
              },
              {
                "color": "dark-blue",
                "value": 0.9
              },
              {
                "color": "super-light-orange",
                "value": 1.2
              },
              {
                "color": "dark-red",
                "value": 2.1
              }
            ]
          },
          "unit": "percentunit"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 19,
        "w": 3,
        "x": 0,
        "y": 26
      },
      "id": 54,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "vertical",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "/.*/",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": false,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "avg(sum by(job) (rate(otelcol_exporter_sent_metric_points_total{}[$__range])))",
          "format": "time_series",
          "hide": true,
          "instant": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "export"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "avg(sum by(job) (rate(otelcol_receiver_accepted_metric_points_total{}[$__range])))",
          "format": "time_series",
          "hide": true,
          "legendFormat": "__auto",
          "range": true,
          "refId": "acc"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "( avg(sum by(job) (rate(otelcol_exporter_sent_metric_points_total{}[$__range]))) /avg(sum by(job) (rate(otelcol_receiver_accepted_metric_points_total{}[$__range]))))",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Export Ratio",
      "transformations": [
        {
          "id": "calculateField",
          "options": {
            "alias": "percent",
            "binary": {
              "left": "avg(sum by(job) (rate(otelcol_exporter_sent_metric_points{}[3600s])))",
              "operator": "/",
              "reducer": "sum",
              "right": "avg(sum by(job) (rate(otelcol_receiver_accepted_metric_points{}[3600s])))"
            },
            "mode": "binary",
            "reduce": {
              "reducer": "sum"
            }
          }
        },
        {
          "id": "organize",
          "options": {
            "excludeByName": {
              "(sum(rate(otelcol_exporter_sent_metric_points{exporter=\"prometheus\"}[1m0s])) )": true,
              "Time": true,
              "avg(sum by(job) (rate(otelcol_exporter_sent_metric_points{}[3600s])))": true,
              "avg(sum by(job) (rate(otelcol_receiver_accepted_metric_points{}[3600s])))": true,
              "{instance=\"otelcol:9464\", job=\"otel\"}": true
            },
            "indexByName": {},
            "renameByName": {
              "Time": "",
              "percent": "Percent",
              "{exporter=\"debug\", instance=\"otelcol:8888\", job=\"otel-collector\", service_instance_id=\"fbfa720a-ebf9-45c8-a79a-9d3b6021a663\", service_name=\"otelcol-contrib\", service_version=\"0.70.0\"}": ""
            }
          }
        }
      ],
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Metrics Signalling Pipelines",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "gridPos": {
        "h": 11,
        "w": 21,
        "x": 3,
        "y": 26
      },
      "id": 25,
      "options": {
        "edges": {},
        "nodes": {
          "mainStatUnit": "flops"
        }
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "\nlabel_join(label_join(\n(rate(otelcol_receiver_accepted_metric_points_total{}[$__interval]))\n, \"id\", \"\", \"transport\", \"receiver\")\n, \"title\", \"\", \"transport\", \"receiver\")\n\nor\n\nlabel_replace(label_replace(\nsum by(service_name) (rate(otelcol_receiver_accepted_spans_total{}[$__interval]))\n, \"id\", \"processor\", \"dummynode\", \"\")\n, \"title\", \"processor\", \"dummynode\", \"\")\n\n\n\nor\nlabel_replace(label_replace(\n(rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n, \"id\", \"$0\", \"processor\", \".*\")\n, \"title\", \"$0\", \"processor\", \".*\")\n\n\n\n\n\nor\nlabel_replace(label_replace(\nsum (rate(otelcol_exporter_sent_metric_points_total{}[$__interval]))\n, \"id\", \"exporter\", \"dummynode\", \"\")\n, \"title\", \"exporter\", \"dummynode\", \"\")\n\nor\nlabel_replace(label_replace(\nsum by(exporter) (rate(otelcol_exporter_sent_metric_points_total{}[$__interval]))\n, \"id\", \"$0\", \"exporter\", \".*\")\n, \"title\", \"$0\", \"exporter\", \".*\")",
          "format": "table",
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "nodes"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "label_join(\nlabel_replace(label_join(\n(rate(otelcol_receiver_accepted_metric_points_total{}[$__interval]))\n\n,\"source\",\"\",\"transport\",\"receiver\")\n,\"target\",\"processor\",\"\",\"\")\n,\"id\",\"-\",\"source\",\"target\")\n\n\nor\n\nlabel_join(\nlabel_replace(label_replace(\n(rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n,\"source\",\"processor\",\"\",\"\")\n,\"target\",\"$0\",\"processor\",\".*\")\n,\"id\",\"-\",\"source\",\"target\")\n\n\n\n\n\nor\n\n\nlabel_join(\nlabel_replace(label_replace(\n(rate(otelcol_processor_batch_batch_send_size_count{}[$__interval]))\n,\"source\",\"$0\",\"processor\",\".*\")\n,\"target\",\"exporter\",\"\",\"\")\n,\"id\",\"-\",\"source\",\"target\")\n\nor\nlabel_join(\nlabel_replace(label_replace(\n(rate(otelcol_exporter_sent_metric_points_total{}[$__interval]))\n,\"source\",\"exporter\",\"\",\"\")\n,\"target\",\"$0\",\"exporter\",\".*\")\n,\"id\",\"-\",\"source\",\"target\")",
          "format": "table",
          "hide": false,
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "edges"
        }
      ],
      "title": "",
      "type": "nodeGraph"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "otelcol_receiver_accepted_metric_points",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "noValue": "no data",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 3,
        "y": 37
      },
      "id": 26,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_accepted_metric_points_total{}[$__rate_interval])) by (receiver,transport)",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Accepted",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "otelcol_receiver_accepted_metric_points\nTotal Accepted ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 3,
        "x": 8,
        "y": 37
      },
      "id": 27,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_accepted_metric_points_total{}[$__rate_interval])) ",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "#EAB839",
                "value": 1
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 11,
        "y": 37
      },
      "id": 28,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "text": {},
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_processor_batch_batch_send_size_sum{}[$__rate_interval]))  by (processor)",
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Batch",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Export ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 3,
        "x": 16,
        "y": 37
      },
      "id": 29,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_sent_metric_points_total{}[$__rate_interval])) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Sent by Exporter",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 5,
        "w": 5,
        "x": 19,
        "y": 37
      },
      "id": 16,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_sent_metric_points_total{}[$__rate_interval])) by (exporter) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Sent",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "noValue": "no data",
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 5,
        "x": 3,
        "y": 42
      },
      "id": 47,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_refused_metric_points_total{}[$__rate_interval])) by (receiver,transport)",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Refused",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Refused \nsum(rate(otelcol_receiver_refused_metric_points{}[$__rate_interval])) ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 8,
        "y": 42
      },
      "id": 48,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "sum(rate(otelcol_receiver_refused_metric_points_total{}[$__rate_interval])) ",
          "legendFormat": "{{receiver}}-{{transport}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total ",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Total Failed Export ",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 3,
        "x": 16,
        "y": 42
      },
      "id": 49,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "auto",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_send_failed_metric_points_total{}[$__rate_interval])) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Total",
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "Sent by Exporter\notelcol_exporter_send_failed_metric_points",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "text",
                "value": null
              }
            ]
          },
          "unit": "none"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 3,
        "w": 5,
        "x": 19,
        "y": 42
      },
      "id": 50,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "auto",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum(rate(otelcol_exporter_send_failed_metric_points_total{}[$__rate_interval])) by (exporter)",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{processor}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Failed",
      "type": "stat"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 45
      },
      "id": 35,
      "panels": [],
      "title": "Prometheus Scrape",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "otelcol prometheus exporter 8888 export rate versus prometheus scrape metrics",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "mappings": [],
          "min": 0,
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "light-blue",
                "value": null
              },
              {
                "color": "semi-dark-red",
                "value": 0
              },
              {
                "color": "super-light-orange",
                "value": 0.4
              },
              {
                "color": "dark-blue",
                "value": 0.9
              },
              {
                "color": "super-light-orange",
                "value": 1.2
              },
              {
                "color": "dark-red",
                "value": 2.1
              }
            ]
          },
          "unit": "percentunit"
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 3,
        "x": 0,
        "y": 46
      },
      "id": 53,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "vertical",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "/.*/",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": false,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "(sum_over_time(scrape_samples_scraped{job=\"otel-collector\"}[$__range])/ count_over_time(scrape_samples_scraped{job=\"otel-collector\"}[$__range])/(5*30)) ",
          "format": "time_series",
          "instant": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "accepted"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "expr": "(sum(rate(otelcol_exporter_sent_metric_points_total{exporter=\"prometheus\"}[$__rate_interval])) )",
          "format": "time_series",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Exported/Scraped",
      "transformations": [
        {
          "id": "calculateField",
          "options": {
            "alias": "percent",
            "binary": {
              "left": "{instance=\"otelcol:9464\", job=\"otel-collector\"}",
              "operator": "/",
              "reducer": "sum",
              "right": "(sum(rate(otelcol_exporter_sent_metric_points{exporter=\"prometheus\"}[1m0s])) )"
            },
            "mode": "binary",
            "reduce": {
              "reducer": "sum"
            }
          }
        },
        {
          "id": "organize",
          "options": {
            "excludeByName": {
              "(sum(rate(otelcol_exporter_sent_metric_points{exporter=\"prometheus\"}[1m0s])) )": true,
              "Time": true,
              "{instance=\"otelcol:9464\", job=\"otel\"}": true
            },
            "indexByName": {},
            "renameByName": {
              "percent": "Percent"
            }
          }
        }
      ],
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "sum_over_time(scrape_samples_scraped[$__range])/ count_over_time(scrape_samples_scraped[$__range])",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 5,
        "x": 3,
        "y": 46
      },
      "id": 37,
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "justifyMode": "auto",
        "orientation": "horizontal",
        "percentChangeColorMode": "standard",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showPercentChange": false,
        "textMode": "value_and_name",
        "wideLayout": true
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum_over_time(scrape_samples_scraped[$__range])/ count_over_time(scrape_samples_scraped[$__range])/(5*30)",
          "format": "time_series",
          "instant": false,
          "legendFormat": "{{job}}/{{instance}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Samples Scraped",
      "type": "stat"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "description": "scrape_samples_scraped{job!=\"\"}\nTotal Samples Scraped",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "continuous-BlYlRd"
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 3,
        "x": 8,
        "y": 46
      },
      "id": 42,
      "options": {
        "minVizHeight": 75,
        "minVizWidth": 75,
        "orientation": "horizontal",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "showThresholdLabels": false,
        "showThresholdMarkers": true,
        "sizing": "auto"
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "sum_over_time(scrape_samples_scraped[$__range])/ count_over_time(scrape_samples_scraped[$__range])/(5*30)",
          "format": "time_series",
          "hide": false,
          "instant": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "B"
        }
      ],
      "title": "Total",
      "transformations": [
        {
          "id": "calculateField",
          "options": {
            "mode": "reduceRow",
            "reduce": {
              "include": [
                "{instance=\"otelcol:8888\", job=\"otel-collector\"}"
              ],
              "reducer": "sum"
            },
            "replaceFields": true
          }
        }
      ],
      "type": "gauge"
    },
    {
      "datasource": {
        "type": "prometheus",
        "uid": "webstore-metrics"
      },
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 8,
        "x": 11,
        "y": 46
      },
      "id": 41,
      "options": {
        "edges": {},
        "nodes": {}
      },
      "pluginVersion": "11.3.0",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "label_replace(label_replace(label_replace(\nsum (scrape_samples_scraped{job!=\"\"}) by (instance)\n, \"id\", \"$0\", \"instance\", \".*\")\n, \"title\", \"$0\", \"instance\", \".*\")\n,\"mainstat\",\"\",\"\",\"\")\n\nor \n\nlabel_replace(label_replace(label_replace(\nsum (scrape_samples_scraped{job!=\"\"})\n, \"id\", \"prometheus\", \"\", \"\")\n, \"title\", \"prometheus\", \"\", \"\")\n,\"mainstat\",\"\",\"\",\"\")\n",
          "format": "table",
          "hide": false,
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "nodes"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "webstore-metrics"
          },
          "editorMode": "code",
          "exemplar": false,
          "expr": "label_join(\nlabel_replace(label_replace(\nsum (scrape_samples_scraped{job!=\"\"}) by (instance)\n,\"source\",\"$0\",\"instance\",\".*\")\n,\"target\",\"prometheus\",\"\",\"\")\n,\"id\",\"-\",\"source\",\"target\")",
          "format": "table",
          "hide": false,
          "instant": true,
          "legendFormat": "__auto",
          "range": false,
          "refId": "edges"
        }
      ],
      "title": "",
      "transformations": [
        {
          "id": "organize",
          "options": {
            "excludeByName": {
              "Time": true
            },
            "indexByName": {},
            "renameByName": {}
          }
        }
      ],
      "type": "nodeGraph"
    },
    {
      "description": "Sent by Exporter",
      "fieldConfig": {
        "defaults": {},
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 5,
        "x": 19,
        "y": 46
      },
      "id": 52,
      "options": {
        "code": {
          "language": "plaintext",
          "showLineNumbers": false,
          "showMiniMap": false
        },
        "content": "\n \n## Prometheus Config\n\n`evaluation_interval:` 30s\n\n`scrape_interval:` 5s",
        "mode": "markdown"
      },
      "pluginVersion": "11.3.0",
      "title": "",
      "type": "text"
    }
  ],
  "preload": false,
  "refresh": false,
  "schemaVersion": 40,
  "tags": [],
  "templating": {
    "list": [
      {
        "allValue": ".*",
        "current": {
          "text": "0.70.0",
          "value": "0.70.0"
        },
        "datasource": {
          "type": "prometheus",
          "uid": "webstore-metrics"
        },
        "definition": "query_result(sum(otelcol_process_uptime{}) by (service_version))\n",
        "hide": 2,
        "includeAll": false,
        "label": "service_version",
        "multi": true,
        "name": "service_version",
        "options": [],
        "query": {
          "query": "query_result(sum(otelcol_process_uptime{}) by (service_version))\n",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "/.*service_version=\"(.*)\".*/",
        "type": "query"
      }
    ]
  },
  "time": {
    "from": "now-15m",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "Opentelemetry Collector Data Flow",
  "uid": "rl5_tea4k",
  "version": 1,
  "weekStart": ""
}
