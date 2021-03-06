import { Box, Divider, Stack, Typography } from "@mui/material";

import { motion } from "framer-motion";
import { Node, NODE_METRIC_COLORS } from "../NodeMap/Node";
import { calculateProgress } from "../NodeMap/Hexagon";

interface BarProgressProps {
  name: string;
  progress: number;
  fill: string;
  width?: number | string;
}

const drawProgress = (progress: number) => {
  return {
    hidden: { width: 0 },
    visible: (i: number) => {
      return {
        width: `${progress}%`,
      };
    },
  };
};

export const BarProgress = ({
  progress,
  fill,
  name,
  width,
}: BarProgressProps) => {
  return (
    <Stack
      direction="row"
      spacing={2}
      sx={{ display: "flex", alignItems: "center" }}
    >
      <Typography variant="subtitle2">{name}</Typography>
      <Box
        sx={{
          width,
          height: "10px",
          backgroundColor: "action.disabledBackground",
          margin: "5px 0",
        }}
      >
        <motion.div
          initial="hidden"
          animate="visible"
          variants={drawProgress(progress)}
          style={{
            width: `${progress}%`,
            height: "10px",
            backgroundColor: fill,
          }}
        />
      </Box>
    </Stack>
  );
};

export const NodeDetail = ({ node }: { node: Node }) => {
  console.log(node.metrics);
  return (
    <div>
      <Typography variant="h6">Node {node.id}</Typography>
      <Divider />
      <Typography variant="subtitle1">Metrics</Typography>
      {node.metrics.map(({ key, value, max }, i) => (
        <div key={key}>
          <BarProgress
            fill={NODE_METRIC_COLORS[i]}
            name={key}
            progress={calculateProgress(value, max) * 100}
          />
        </div>
      ))}
    </div>
  );
};
