package virtualmemory
import (
	"github.com/kokdot/go-musthave-devops/internal/def"
    "github.com/shirou/gopsutil/v3/mem"
)

func GetData(virtualMemoryMap *def.GaugeMap) *def.GaugeMap {
    v, _ := mem.VirtualMemory()
	(*virtualMemoryMap)["TotalMemory"] = def.Gauge(v.Total)
	(*virtualMemoryMap)["FreeMemory"] = def.Gauge(v.Free)
	(*virtualMemoryMap)["CPUutilization1"] = def.Gauge(v.UsedPercent)
	return virtualMemoryMap
}
