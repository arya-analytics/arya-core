import { Page, PageHeading } from "../Page/Page";
import { NodeMap } from "../NodeMap/NodeMap";
import { GrainOutlined } from "@mui/icons-material";
import { Node } from "../NodeMap/Node";
import { useState } from "react";
import { ClusterInfo } from "./ClusterInfo";

const DUMMY_NODES: Node[] = [
  {
    id: 1,
    metrics: [
      {
        key: "Memory",
        value: 12,
        max: 32,
      },
      {
        key: "CPU",
        value: 4.32,
        max: 8.0,
      },
      {
        key: "Storage",
        value: 300,
        max: 900,
      },
    ],
  },
  {
    id: 2,
    metrics: [
      {
        key: "Memory",
        value: 12,
        max: 16,
      },
      {
        key: "CPU",
        value: 7.32,
        max: 8.0,
      },
      {
        key: "Storage",
        value: 450,
        max: 900,
      },
    ],
  },
  {
    id: 3,
    metrics: [
      {
        key: "Memory",
        value: 96,
        max: 128,
      },
      {
        key: "CPU",
        value: 4.32,
        max: 12.0,
      },
      {
        key: "Storage",
        value: 722,
        max: 900,
      },
      {
        key: "Active Channels",
        value: 815,
        max: 902,
      },
      {
        key: "Write Throughput",
        value: 300000,
        max: 900000,
      },
    ],
  },
];

export const Cluster = () => {
  const [selectedNode, setSelectedNode] = useState<Node | null>(DUMMY_NODES[0]);
  return (
    <Page direction="column">
      <PageHeading variant="heading" icon={<GrainOutlined />}>
        Cluster
      </PageHeading>
      <Page direction="row">
        <NodeMap nodes={DUMMY_NODES} onSelect={setSelectedNode} />
        <ClusterInfo selectedNode={selectedNode} />
      </Page>
    </Page>
  );
};
