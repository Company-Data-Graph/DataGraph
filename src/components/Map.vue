<template>
  <div id="graph">
    <div :class="[isLoading ? 'd-flex' : 'd-none']" class="graph__loader align-items-center justify-content-center">
      <div>Загружено пакетов: {{this.countOfLoadedPackages}} / {{this.countOfPackages}} </div>
    </div>
    <div :class="[!isLoading && isDataLoading ? 'd-flex' : 'd-none']"
      class="graph__loader graph__loader_t align-items-center justify-content-center">
    </div>
    <div :class="[isLoadingPackagesFailure ? 'd-flex' : 'd-none']" class="graph__loader align-items-center justify-content-center"><div> Error </div></div>
    <div id="graph_container"
         @pointermove="graphInteractionMouseMove"
         @wheel.stop.prevent="graphInteractionMouseWheel"
         @pointerdown="graphInteractionMouseDown"
         @pointerup="graphInteractionMouseUp"
         ref="graph_container">
      <canvas id="graph_canvas" ref="graph_canvas"></canvas>
      </div>
  </div>
</template>

<script>
import * as PIXI from 'pixi.js'
import axios from 'axios'

export default {
  name: "InteractiveMap",
  data() {
    return {
      config: {
        container: "graph_container",
        canvas_id: "graph_canvas",
        background_color: 0xE1E1E1,
        min_zoom: 0.2,
        max_zoom: 2.,
      },
      nodes: [],
      links: [],
      pixiApp: {
        'app': {type: Object},
        'stages': {type: Object},
        'graphics': []
      },
      hoveredGraphics: undefined,
      interactionParams: {
        delY: 0, ofX: 0, ofY: 0, offsetX: 0, offsetY: 0,
        lastPos: null, isMoved: false, isNotClicked: false, isOnDown: false
      },
      recalculateGraphZoom: 0.4,
      graphMode: undefined,
      isLoadingPackages: false,
      isLoadingPackagesFailure: false,
      tooltipParams: {
        isShow: false,
        transform: {x: 1000, y: 200},
        offset: {x: 1000, y: 200 },
        isArrowStart: true,
        data: {
          title: "",
        },
      },
      notActiveColor: 0xD9D9D9,
      selectColor: 0xFF00C7
    }
  },
  props: {
    graphInteractions: {
      required: false
    },
    isDataLoading: {
      required: true,
      default: false
    },
    currentFiltersCompany: {
      required: false,
      default: undefined
    },
    currentFiltersProduct: {
      required: false,
      default: undefined
    },
    setPanelState: {
      type: Function,
      required: true
    },
    selectedNode: {
      required: false,
      default: undefined
    }
  },
  watch: {
    graphInteractions: {
      handler(newValue) {
        this.zoom(newValue.oldScaleFactor - newValue.newScaleFactor, this.$refs.graph_container.offsetWidth / 2, this.$refs.graph_container.offsetHeight / 2)
      },
      deep: true
    },
    currentFiltersCompany: {
      handler(newValue) {
        this.recalculateNodes()
      },
      deep: true
    },
    currentFiltersProduct: {
      handler(newValue) {
        this.recalculateNodes()
      },
      deep: true
    },
    selectedNode: {
      handler(newValue) {
        console.log("new select 1")
        this.recalculateNodes()
      },
      deep: true
    }
  },
  methods: {
    recalculateNodes: function () {
      this.clearStages([this.pixiApp.stages.mainStage])
      this.initNodes(this.nodes, this.links)
    },
    clearStages: function (stages) {
      stages.forEach((stage) => {
        for (let i = stage.children.length - 1; i >= 0; i--) {
          stage.removeChild(stage.children[i])
        }}
      )
    },
    initPixiApplication: function (config) {
      let canvas = this.$refs.graph_canvas;
      let app = new PIXI.Application({
        width: this.$refs.graph_container.offsetWidth,
        height: this.$refs.graph_container.offsetHeight,
        backgroundColor: config.background_color,
        forceCanvas: true, view: canvas
      })
      this.$refs.graph_container.appendChild(app.view);
      let mainStage = new PIXI.Container();
      app.stage.addChild(mainStage)
      let stages = {'mainStage': mainStage}
      return app, stages
    },
    nodeIdInFilter: function(id) {
      console.log("test")
      if (this.nodes.filter(el => el.id == id)[0].nodeType == "Компания") {
        if (this.currentFiltersCompany == undefined) {
          return true
        }
        return this.currentFiltersCompany.filter(el => el == id).length != 0
      } else {
        if (this.currentFiltersProduct == undefined) {
          return true
        }
        return this.currentFiltersProduct.filter(el => el == id).length != 0
      }


    },
    initNodes: function (nodes, links) {
      for (let id = 0; id < links.length; id++) {
        if (links[id] != null) {
          let current = links[id]
          let graphics = new PIXI.Graphics();
          let color = PIXI.utils.string2hex("0x" + nodes.filter(el => el.id == current.source)[0].color.substring(1).trim())
          if (this.nodeIdInFilter(current.source) == false || this.nodeIdInFilter(current.target) == false) {
            color = this.notActiveColor
          } else {
            graphics.interactive = true
          }
          graphics.beginFill(color, 1);
          let angle = this.drawLine(graphics, nodes.filter(el => el.id == current.source)[0], nodes.filter(el => el.id == current.target)[0], color)
          graphics.closePath();
          graphics.endFill();
          let sourceNode = nodes.filter(el => el.id == current.source)[0]
          graphics.position.set( sourceNode.x, sourceNode.y );
          graphics.pivot.x = sourceNode.x
          graphics.pivot.y = sourceNode.y
          graphics.closePath();
          graphics.endFill();
          graphics.rotation = angle
          if (this.nodeIdInFilter(current.source) == true  || this.nodeIdInFilter(current.target) == false) {
            graphics.on('pointerdown', () => {
              if (this.interactionParams.isMoved || this.interactionParams.isNotClicked)
                return;
              this.setPanelState(true)
              let sourceNodeType = this.nodes.filter(el => el.id == current.source)[0].nodeType
              this.$router.push({path: '/line',
                query: {
                  mode: sourceNodeType == "Компания" ? "company" : "product",
                  source: current.source,
                  target: current.target
                }
              });
            })
          }
          graphics.id = current.id
          graphics.nodeType = current.nodeType
          this.pixiApp.graphics.push(graphics)
          this.pixiApp.stages.mainStage.addChild(graphics);
        }
      }
      for (let id = 0; id < nodes.length; id++) {
        if (nodes[id] != null) {
          let current = nodes[id]
          let graphics = new PIXI.Graphics();
          let color = PIXI.utils.string2hex("0x" + current.color.substring(1).trim())
          if (this.nodeIdInFilter(current.id) == false) {
            color = this.notActiveColor
          } else {
            graphics.interactive = true
            graphics.on('pointerdown', () => {
              if (this.interactionParams.isMoved || this.interactionParams.isNotClicked)
                return;
              this.setPanelState(true)
              if (graphics.nodeType == "Продукт")
                this.$router.push({path: '/product/' + graphics.id});
              if (graphics.nodeType == "Компания")
                this.$router.push({path: '/company/' + graphics.id});
            })
          }
          console.log(current.svg)
          graphics.beginFill(color, 1);
          this.drawNode(graphics, current)
          graphics.closePath();
          graphics.endFill();
          graphics.id = current.id
          graphics.nodeType = current.nodeType
          this.pixiApp.graphics.push(graphics)
          this.pixiApp.stages.mainStage.addChild(graphics);
        }
      }
      console.log(links)
    },
    drawNode: function (graphics, node) {
      if (this.selectedNode != undefined && this.selectedNode.nodeType == "node" &&
      this.selectedNode.id == node.id) {
        graphics.lineStyle(10, this.selectColor, 1);
      } else {
        graphics.lineStyle(3, node.color, 1);
      }
      graphics.drawCircle(node.x, node.y, (node.nodeType == "Продукт") ? 50 : 100)
      const style = new PIXI.TextStyle({
        fontFamily: 'Arial',
        fontSize: (node.nodeType == "Продукт") ? 25 : 50,
        fontWeight: 'bold',
      });
      const headerLetter = new PIXI.Text(node.name[0], style)
      headerLetter.x = node.x - style.fontSize / 2
      headerLetter.y = node.y - style.fontSize / 2
      graphics.addChild(headerLetter)
    },

    drawLine: function (graphics, source, target, color) {
      let line = new PIXI.Graphics()
      if (this.selectedNode != undefined && this.selectedNode.nodeType == "line" &&
      this.selectedNode.source == source.id && this.selectedNode.target == target.id) {
        color = this.selectColor;
        console.log("get")
      } else {
        line.lineStyle(3, color, 1);
      }
      const lineSize = 5;

      let magnitude = Math.sqrt(Math.pow(target.x - source.x, 2) + Math.pow(target.y - source.y, 2))

      let angle = Math.atan( Math.abs(target.x - source.x) / Math.abs(target.y - source.y))
      if (target.x > source.x) {
        angle = -angle;
      }

      var points = [
          { x: (source.x + lineSize), y: source.y},
          { x: (source.x - lineSize), y: source.y},
          { x: (source.x - lineSize), y: source.y + magnitude},
          { x: (source.x + lineSize), y: source.y + magnitude},
      ]


      graphics.beginFill(color, 1);
      graphics.moveTo(points[0].x, points[0].y)
      graphics.lineTo(points[1].x, points[1].y)
      graphics.lineTo(points[2].x, points[2].y)
      graphics.lineTo(points[3].x, points[3].y)
      graphics.lineTo(points[0].x, points[0].y)
      return angle

    },
    zoom: function (s, x, y) {
      s = -s
      s = s > 0 ? 2 : 0.5;
      let stages = this.pixiApp.stages
      let worldPos = {
        x: (x - stages.mainStage.x) / stages.mainStage.scale.x,
        y: (y - stages.mainStage.y) / stages.mainStage.scale.y
      };
      let newScale = {
        x: stages.mainStage.scale.x * s,
        y: stages.mainStage.scale.y * s
      };
      newScale.x = newScale.x >= this.config.max_zoom ? this.config.max_zoom :
          newScale.x <= this.config.min_zoom ? this.config.min_zoom : newScale.x;
      newScale.y = newScale.y >= this.config.max_zoom ? this.config.max_zoom :
          newScale.y <= this.config.min_zoom ? this.config.min_zoom : newScale.y;
      console.log("Before zoom: " + this.graphMode + " " + newScale.x)
      if (newScale.x >= this.recalculateGraphZoom && this.graphMode == "short") {
        this.fetchingNodes("full")
      }
      if (newScale.x < this.recalculateGraphZoom && this.graphMode == "full") {
        this.fetchingNodes("short")
      }

      let newScreenPos = {
        x: (worldPos.x) * newScale.x + stages.mainStage.x,
        y: (worldPos.y) * newScale.y + stages.mainStage.y
      };

      this.pixiApp.stages.mainStage.x -= (newScreenPos.x - x);
      this.pixiApp.stages.mainStage.y -= (newScreenPos.y - y);
      this.pixiApp.stages.mainStage.scale.x = newScale.x;
      this.pixiApp.stages.mainStage.scale.y = newScale.y;
    },
    graphInteractionMouseWheel: function (e) {
      this.interactionParams.delY += e.deltaY
      this.interactionParams.ofX += e.offsetX
      this.interactionParams.ofY += e.offsetY
      if (window.wheelTimeout) {
        window.clearTimeout(window.wheelTimeout);
      }
      this.zoom(e.deltaY, e.offsetX, e.offsetY)
    },
    graphInteractionMouseDown: function (e) {
      this.interactionParams.isOnDown = true;
      this.interactionParams.lastPos = {x: e.offsetX, y: e.offsetY};
      this.interactionParams.offsetX = 0, this.interactionParams.offsetY = 0
    },
    graphInteractionMouseUp: function (event) {
      this.interactionParams.isOnDown = false;
      if (this.interactionParams.isMoved) {
        this.interactionParams.isMoved = false;
        this.interactionParams.isNotClicked = true;
        setTimeout(() => {
          this.interactionParams.isNotClicked = false;
        }, 10);
      } else {
        this.interactionParams.isOnDown = false;
        this.interactionParams.isMoved = false;
        this.interactionParams.isNotClicked = false;
      }
    },
    graphInteractionMouseMove: function (e) {
      if (this.interactionParams.isOnDown) {
        this.interactionParams.isMoved = true;
      }
      if (this.interactionParams.lastPos) {
        this.interactionParams.offsetX += (e.offsetX - this.interactionParams.lastPos.x);
        this.interactionParams.offsetY += (e.offsetY - this.interactionParams.lastPos.y);
        this.interactionParams.lastPos = {x: e.offsetX, y: e.offsetY};
        if (this.interactionParams.isMoved) {
          this.pixiApp.stages.mainStage.position.x += this.interactionParams.offsetX
          this.pixiApp.stages.mainStage.position.y += this.interactionParams.offsetY
          this.interactionParams.offsetX = 0, this.interactionParams.offsetY = 0
        }
      }
      this.tooltipParams.offset = {x: e.offsetX, y: e.offsetY}
    },
    fetchingNodes: function(mode) {
      const k = 2;
      this.isLoadingPackages = true
      this.isLoadingPackagesFailure = false
      this.graphMode = mode
      axios.get("http://141.95.127.215:7328/get:"+mode).then(response => {
        this.nodes = response.data.nodes
        this.links = response.data.links
        this.nodes.map(node => {
          node.x = node.x * k
          node.y = node.y * k
        })
        this.recalculateNodes()
        this.isLoadingPackages = false
        console.log("After zoom: " + this.graphMode)
      }).catch(e => {
        console.log(e)
        this.isLoadingPackagesFailure = true
        this.isLoadingPackages = false
      })

    }
  },
  mounted() {
    this.pixiApp.app, this.pixiApp.stages = this.initPixiApplication(this.config)
    this.fetchingNodes("short")
  },
}

</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
#graph {
  height: 100vh;
  width: 100%;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 1;
}
.graph__loader {
  position: absolute;
  height: 100%;
  width: 75%;
  top: 0;
  right: 0;
}

.graph__loader {
  background: white;
}

.graph__loader_t {
  background: rgba(255,255,255,.6);
}

#graph_container {
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.graph_tooltip {
  position: absolute;
  pointer-events: none;
}

</style>

