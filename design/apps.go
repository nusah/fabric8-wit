package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

// SimpleSpace describe a space
var simpleSpace = a.Type("SimpleSpace", func() {
	a.Description(`a space consisting of multiple applications`)
	a.Attribute("id", d.UUID)
	a.Attribute("name", d.String)
	a.Attribute("applications", a.ArrayOf(simpleApp))
	a.Required("applications")
})

// SimpleApp describe an application within a space
var simpleApp = a.Type("SimpleApp", func() {
	a.Description(`a description of an application`)
	a.Attribute("id", d.UUID)
	a.Attribute("name", d.String)
	a.Attribute("pipeline", a.ArrayOf(simpleDeployment))
	a.Required("pipeline")
})

// simpleDeployment describe an element of an application pipeline
var simpleDeployment = a.Type("SimpleDeployment", func() {
	a.Description(`a deployment (a step in a pipeline, e.g. 'build')`)
	a.Attribute("id", d.UUID)
	a.Attribute("name", d.String)
	a.Attribute("version", d.String)
	a.Attribute("pods", podStats)
})

// simpleDeployment describe an element of an application pipeline
var simpleEnvironment = a.Type("SimpleEnvironment", func() {
	a.Description(`a shared environment`)
	a.Attribute("id", d.UUID)
	a.Attribute("name", d.String)
	a.Attribute("quota", envStats)
})

var envStats = a.Type("EnvStats", func() {
	a.Description("resource usage and quotas for an environment")
	a.Attribute("cpucores", envStatCores)
	a.Attribute("memory", envStatMemory)
})

var envStatCores = a.Type("EnvStatCores", func() {
	a.Description(`CPU core stats`)
	a.Attribute("used", d.Number)
	a.Attribute("quota", d.Number)
})

var envStatMemory = a.Type("EnvStatMemory", func() {
	a.Description(`memory stats`)
	a.Attribute("used", d.Number)
	a.Attribute("quota", d.Number)
	a.Attribute("units", d.String)
})

var podStats = a.Type("PodStats", func() {
	a.Description(`pod stats`)
	a.Attribute("starting", d.Integer)
	a.Attribute("running", d.Integer)
	a.Attribute("stopping", d.Integer)
	a.Attribute("total", d.Integer)
})

var timedNumberTuple = a.Type("TimedNumberTuple", func() {
	a.Description("a set of time and number values")
	a.Attribute("time", d.Number)
	a.Attribute("value", d.Number)
})

var simpleDeploymentStats = a.Type("SimpleDeploymentStats", func() {
	a.Description("current deployment stats")
	a.Attribute("cores", timedNumberTuple)
	a.Attribute("memory", timedNumberTuple)
})

var simpleDeploymentStatSeries = a.Type("SimpleDeploymentStatSeries", func() {
	a.Description("pod stat series")
	a.Attribute("start", d.Number)
	a.Attribute("end", d.Number)
	a.Attribute("memory", a.ArrayOf(timedNumberTuple))
	a.Attribute("cores", a.ArrayOf(timedNumberTuple))
})

var simpleSpaceSingle = JSONSingle(
	"SimpleSpace", "Holds a single response to a space request",
	simpleSpace,
	nil)

var simpleAppSingle = JSONSingle(
	"SimpleApplication", "Holds a single response to a space/application request",
	simpleApp,
	nil)

var simpleEnvironmentSingle = JSONSingle(
	"SimpleEnvironment", "Holds a single response to a space/environment request",
	simpleEnvironment,
	nil)

var simpleEnvironmentMultiple = JSONList(
	"SimpleEnvironment", "Holds a response to a space/environment request",
	simpleEnvironment,
	nil,
	nil)

var simplePod = a.Type("SimplePod", func() {
	a.Description("wrapper for a kubernetes Pod")
	a.Attribute("pod", d.Any)
})

var simplePodMultiple = JSONList(
	"SimplePod", "Holds a list of pods",
	simplePod,
	nil,
	nil)

var simpleDeploymentSingle = JSONSingle(
	"SimpleDeployment", "Holds a single response to a space/application/deployment request",
	simpleDeployment,
	nil)

var simpleDeploymentStatsSingle = JSONSingle(
	"SimpleDeploymentStats", "Holds a single response to a space/application/deployment/stats request",
	simpleDeploymentStats,
	nil)

var simpleDeploymentStatSeriesSingle = JSONSingle(
	"SimpleDeploymentStatSeries", "HOlds a response to a stat series query",
	simpleDeploymentStatSeries,
	nil)

var simpleEnvironmentStatSingle = JSONSingle(
	"EnvStats", "Holds a single response to a pipeline/stats request",
	envStats,
	nil)

var _ = a.Resource("apps", func() {
	a.BasePath("/apps")

	// An auth token is required to call the auth API to get an OpenShift auth token.
	a.Security("jwt")

	a.Action("showSpace", func() {
		a.Routing(
			a.GET("/spaces/:spaceID"),
		)
		a.Description("list applications in a space")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
		})
		a.Response(d.OK, simpleSpaceSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showSpaceApp", func() {
		a.Routing(
			a.GET("/spaces/:spaceID/applications/:appName"),
		)
		a.Description("list application")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
			a.Param("appName", d.String, "Name of the application")
		})
		a.Response(d.OK, simpleAppSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showSpaceAppDeployment", func() {
		a.Routing(
			a.GET("/spaces/:spaceID/applications/:appName/deployments/:deployName"),
		)
		a.Description("list pipe element")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
			a.Param("appName", d.String, "Name of the application")
			a.Param("deployName", d.String, "Name of the pipe deployment")
		})
		a.Response(d.OK, simpleDeploymentSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showDeploymentStats", func() {
		a.Routing(
			a.GET("/spaces/:spaceID/applications/:appName/deployments/:deployName/stats"),
		)
		a.Description("get deployment statistics")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
			a.Param("appName", d.String, "Name of the application")
			a.Param("deployName", d.String, "Name of the deployment")
			a.Param("start", d.Number, "start time in millis")
		})
		a.Response(d.OK, simpleDeploymentStatsSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showDeploymentStatSeries", func() {
		a.Routing(
			a.GET("/spaces/:spaceID/applications/:appName/deployments/:deployName/statseries"),
		)
		a.Description("list deployment statistics")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
			a.Param("appName", d.String, "Name of the application")
			a.Param("deployName", d.String, "Name of the deployment")
			a.Param("start", d.Number, "start time in millis")
			a.Param("end", d.Number, "end time in millis")
			a.Param("limit", d.Integer, "maximum number of data points to return")
		})
		a.Response(d.OK, simpleDeploymentStatSeriesSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("setDeployment", func() {
		a.Routing(
			a.PUT("/spaces/:spaceID/applications/:appName/deployments/:deployName/control"),
		)
		a.Description("set deployment pod count")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
			a.Param("appName", d.String, "Name of the application")
			a.Param("deployName", d.String, "Name of the deployment")
			a.Param("podCount", d.Integer, "desired running pod count")
		})
		a.Response(d.OK)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showSpaceEnvironments", func() {
		a.Routing(
			a.GET("/spaces/:spaceID/environments"),
		)
		a.Description("list all environments for a space")
		a.Params(func() {
			a.Param("spaceID", d.UUID, "ID of the space")
		})
		a.Response(d.OK, simpleEnvironmentMultiple)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showEnvironment", func() {
		a.Routing(
			a.GET("/environments/:envName"),
		)
		a.Description("list environment")
		a.Params(func() {
			a.Param("envName", d.String, "Name of the environment")
		})
		a.Response(d.OK, simpleEnvironmentSingle)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

	a.Action("showEnvAppPods", func() {
		a.Routing(
			a.GET("/environments/:envName/applications/:appName/pods"),
		)
		a.Description("list application pods")
		a.Params(func() {
			a.Param("envName", d.String, "Name of the environment")
			a.Param("appName", d.String, "Name of the application")
		})
		// TODO - find a way to use predefined structs in goa DSL
		// until then, hand code JSON response here instead of []v1.Pod
		a.Response(d.OK, "application/json")
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
	})

})
