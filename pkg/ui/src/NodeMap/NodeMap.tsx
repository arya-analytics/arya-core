import React from "react";
import { Node, NodePin } from "./Node";
import { NoiseAwareOutlined } from "@mui/icons-material";
import { Page, PageHeading } from "../Page/Page";
import { Canvas } from "../Canvas/Canvas";

interface NodeMapProps {
  onSelect?: (node: Node) => void;
  nodes: Node[];
}

export const NodeMap = ({ nodes, onSelect }: NodeMapProps) => {
  return (
    <Page direction="column" xs={8}>
      <PageHeading variant="subheading" icon={<NoiseAwareOutlined />}>
        Node Map
      </PageHeading>
      <Page direction="row">
        <Canvas>
          {nodes.map((node) => {
            return <NodePin node={node} onClick={onSelect} />;
          })}
        </Canvas>
      </Page>
    </Page>
  );
};
