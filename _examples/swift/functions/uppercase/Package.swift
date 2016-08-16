import PackageDescription

let package = Package(
    name: "main",
    dependencies: [
        .Package(url: "https://github.com/paulofaria/swift-apex.git", majorVersion: 0, minor: 1),
    ]
)
